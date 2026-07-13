package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/mysqlx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type TripRepository struct {
	db *sql.DB
}

type tripPlanPayload struct {
	Summary *trip.TripSummary  `json:"summary"`
	Days    []*trip.TripDay    `json:"days"`
	Budget  []*trip.TripBudget `json:"budget"`
	Alerts  []string           `json:"alerts"`
}

func NewTripRepository(cfg config.MySQLConfig) (*TripRepository, error) {
	db, err := mysqlx.New(cfg)
	if err != nil {
		return nil, err
	}
	repo := &TripRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *TripRepository) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS trips (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) UNSIGNED NOT NULL,
    title VARCHAR(128) NOT NULL DEFAULT '',
    subtitle VARCHAR(255) NOT NULL DEFAULT '',
    destination VARCHAR(128) NOT NULL DEFAULT '',
    date_range VARCHAR(64) NOT NULL DEFAULT '',
    day_count INT NOT NULL DEFAULT 0,
    people VARCHAR(64) NOT NULL DEFAULT '',
    budget_level VARCHAR(64) NOT NULL DEFAULT '',
    source_question TEXT,
    source_reply MEDIUMTEXT,
    plan_json LONGTEXT,
    saved TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_trips_user_created (user_id, created_at),
    KEY idx_trips_user_destination (user_id, destination)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

func (r *TripRepository) Create(ctx context.Context, input *repository.Trip) (*repository.Trip, error) {
	stored := cloneTrip(input)
	if stored.CreatedAt.IsZero() {
		stored.CreatedAt = time.Now()
	}
	stored.UpdatedAt = stored.CreatedAt
	planJSON, err := marshalTripPlan(stored)
	if err != nil {
		return nil, err
	}
	result, err := r.db.ExecContext(ctx, `
INSERT INTO trips(
    user_id, title, subtitle, destination, date_range, day_count, people,
    budget_level, source_question, source_reply, plan_json, saved, created_at, updated_at
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		stored.UserID,
		stored.Title,
		stored.Subtitle,
		stored.Destination,
		stored.DateRange,
		stored.DayCount,
		stored.People,
		stored.BudgetLevel,
		stored.SourceQuestion,
		stored.SourceReply,
		planJSON,
		boolToInt(stored.Saved),
		stored.CreatedAt,
		stored.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	stored.ID = strconv.FormatInt(id, 10)
	return stored, nil
}

func (r *TripRepository) ListByUser(ctx context.Context, userID int64, limit int) ([]*repository.Trip, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, title, subtitle, destination, date_range, day_count, people,
       budget_level, source_question, source_reply, plan_json, saved, created_at, updated_at, deleted_at
FROM trips
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*repository.Trip, 0)
	for rows.Next() {
		item, err := scanTrip(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *TripRepository) LatestByUser(ctx context.Context, userID int64) (*repository.Trip, bool, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_id, title, subtitle, destination, date_range, day_count, people,
       budget_level, source_question, source_reply, plan_json, saved, created_at, updated_at, deleted_at
FROM trips
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT 1`, userID)
	return scanTripRow(row)
}

func (r *TripRepository) FindByIDForUser(ctx context.Context, userID int64, tripID string) (*repository.Trip, bool, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_id, title, subtitle, destination, date_range, day_count, people,
       budget_level, source_question, source_reply, plan_json, saved, created_at, updated_at, deleted_at
FROM trips
WHERE id = ? AND user_id = ? AND deleted_at IS NULL
LIMIT 1`, tripID, userID)
	return scanTripRow(row)
}

func (r *TripRepository) DeleteByIDForUser(ctx context.Context, userID int64, tripID string) bool {
	now := time.Now()
	result, err := r.db.ExecContext(ctx, `
UPDATE trips
SET deleted_at = ?, updated_at = ?
WHERE id = ? AND user_id = ? AND deleted_at IS NULL`, now, now, tripID, userID)
	if err != nil {
		return false
	}
	affected, err := result.RowsAffected()
	return err == nil && affected > 0
}

func (r *TripRepository) Update(ctx context.Context, userID int64, tripID string, apply func(*repository.Trip)) (*repository.Trip, bool, error) {
	item, found, err := r.FindByIDForUser(ctx, userID, tripID)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	apply(item)
	now := time.Now()
	item.UpdatedAt = now
	planJSON, err := marshalTripPlan(item)
	if err != nil {
		return nil, false, err
	}
	_, err = r.db.ExecContext(ctx, `
UPDATE trips
SET title = ?, subtitle = ?, destination = ?, date_range = ?, day_count = ?,
    people = ?, budget_level = ?, plan_json = ?, updated_at = ?
WHERE id = ? AND user_id = ? AND deleted_at IS NULL`,
		item.Title,
		item.Subtitle,
		item.Destination,
		item.DateRange,
		item.DayCount,
		item.People,
		item.BudgetLevel,
		planJSON,
		item.UpdatedAt,
		tripID,
		userID,
	)
	if err != nil {
		return nil, false, err
	}
	return item, true, nil
}

type tripScanner interface {
	Scan(dest ...any) error
}

func scanTripRow(row tripScanner) (*repository.Trip, bool, error) {
	item, err := scanTrip(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return item, true, nil
}

func scanTrip(row tripScanner) (*repository.Trip, error) {
	var item repository.Trip
	var id int64
	var planJSON sql.NullString
	var saved int
	var deletedAt sql.NullTime
	if err := row.Scan(
		&id,
		&item.UserID,
		&item.Title,
		&item.Subtitle,
		&item.Destination,
		&item.DateRange,
		&item.DayCount,
		&item.People,
		&item.BudgetLevel,
		&item.SourceQuestion,
		&item.SourceReply,
		&planJSON,
		&saved,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	); err != nil {
		return nil, err
	}
	item.ID = strconv.FormatInt(id, 10)
	item.Saved = saved == 1
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}
	if planJSON.Valid {
		applyTripPlanJSON(&item, planJSON.String)
	}
	return &item, nil
}

func marshalTripPlan(input *repository.Trip) (string, error) {
	payload := tripPlanPayload{
		Summary: cloneTripSummary(input.Summary),
		Days:    cloneTripDays(input.Days),
		Budget:  cloneTripBudget(input.Budget),
		Alerts:  append([]string(nil), input.Alerts...),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func applyTripPlanJSON(input *repository.Trip, raw string) {
	var payload tripPlanPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return
	}
	input.Summary = cloneTripSummary(payload.Summary)
	input.Days = cloneTripDays(payload.Days)
	input.Budget = cloneTripBudget(payload.Budget)
	input.Alerts = append([]string(nil), payload.Alerts...)
}

func cloneTrip(input *repository.Trip) *repository.Trip {
	if input == nil {
		return nil
	}
	copied := *input
	copied.Summary = cloneTripSummary(input.Summary)
	copied.Days = cloneTripDays(input.Days)
	copied.Budget = cloneTripBudget(input.Budget)
	copied.Alerts = append([]string(nil), input.Alerts...)
	if input.DeletedAt != nil {
		deletedAt := *input.DeletedAt
		copied.DeletedAt = &deletedAt
	}
	return &copied
}

func cloneTripSummary(input *trip.TripSummary) *trip.TripSummary {
	if input == nil {
		return nil
	}
	return &trip.TripSummary{
		Date:   cloneStringPtr(input.Date),
		Days:   cloneStringPtr(input.Days),
		People: cloneStringPtr(input.People),
		Budget: cloneStringPtr(input.Budget),
	}
}

func cloneTripDays(input []*trip.TripDay) []*trip.TripDay {
	result := make([]*trip.TripDay, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripDay{
			Day:     item.Day,
			Title:   cloneStringPtr(item.Title),
			Route:   cloneStringPtr(item.Route),
			Food:    cloneStringPtr(item.Food),
			Hotel:   cloneStringPtr(item.Hotel),
			Tips:    cloneTripTips(item.Tips),
			Weather: cloneStringPtr(item.Weather),
		})
	}
	return result
}

func cloneTripTips(input []*trip.TripTip) []*trip.TripTip {
	result := make([]*trip.TripTip, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripTip{
			Icon:  cloneStringPtr(item.Icon),
			Title: cloneStringPtr(item.Title),
			Text:  cloneStringPtr(item.Text),
		})
	}
	return result
}

func cloneTripBudget(input []*trip.TripBudget) []*trip.TripBudget {
	result := make([]*trip.TripBudget, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripBudget{
			Label:  cloneStringPtr(item.Label),
			Amount: cloneStringPtr(item.Amount),
		})
	}
	return result
}

func cloneStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := *input
	return &value
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
