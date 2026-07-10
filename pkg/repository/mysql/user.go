package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/mysqlx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(cfg config.MySQLConfig) (*UserRepository, error) {
	db, err := mysqlx.New(cfg)
	if err != nil {
		return nil, err
	}
	repo := &UserRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *UserRepository) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    phone VARCHAR(32) DEFAULT NULL,
    email VARCHAR(128) DEFAULT NULL,
    password_hash VARCHAR(255) NOT NULL DEFAULT '',
    nickname VARCHAR(64) NOT NULL DEFAULT '',
    avatar_url VARCHAR(512) NOT NULL DEFAULT '',
    gender TINYINT NOT NULL DEFAULT 0,
    birthday DATE DEFAULT NULL,
    home_city VARCHAR(128) NOT NULL DEFAULT '',
    current_city VARCHAR(128) NOT NULL DEFAULT '',
    status TINYINT NOT NULL DEFAULT 1,
    last_login_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_phone (phone),
    UNIQUE KEY uk_users_email (email),
    KEY idx_users_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		return err
	}
	if _, err = r.db.ExecContext(ctx, `ALTER TABLE users ADD COLUMN current_city VARCHAR(128) NOT NULL DEFAULT '' AFTER home_city`); err != nil && !isDuplicateColumnError(err, "current_city") {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS user_preferences (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) UNSIGNED NOT NULL,
    preference_value VARCHAR(128) NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_preference (user_id, preference_value),
    KEY idx_user_preferences_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT(20) UNSIGNED NOT NULL,
    trip_reminder_enabled TINYINT NOT NULL DEFAULT 1,
    price_reminder_enabled TINYINT NOT NULL DEFAULT 1,
    personalized_recommend_enabled TINYINT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

func (r *UserRepository) Create(ctx context.Context, input *repository.User) (*repository.User, error) {
	now := time.Now()
	storedUser := cloneUser(input)
	storedUser.CreatedAt = now
	storedUser.UpdatedAt = now
	result, err := r.db.ExecContext(ctx, `
INSERT INTO users(phone, password_hash, nickname, avatar_url, home_city, current_city, status, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		storedUser.Phone,
		storedUser.PasswordHash,
		storedUser.Nickname,
		storedUser.AvatarURL,
		storedUser.HomeCity,
		storedUser.CurrentCity,
		storedUser.Status,
		storedUser.CreatedAt,
		storedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	storedUser.ID = strconv.FormatInt(insertID, 10)
	return storedUser, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*repository.User, bool) {
	return r.findOne(ctx, "phone = ?", phone)
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*repository.User, bool) {
	return r.findOne(ctx, "id = ?", userID)
}

func (r *UserRepository) Update(ctx context.Context, userID string, apply func(*repository.User)) (*repository.User, bool) {
	storedUser, ok := r.FindByID(ctx, userID)
	if !ok {
		return nil, false
	}
	apply(storedUser)
	storedUser.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET nickname = ?, avatar_url = ?, home_city = ?, current_city = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL`,
		storedUser.Nickname,
		storedUser.AvatarURL,
		storedUser.HomeCity,
		storedUser.CurrentCity,
		storedUser.UpdatedAt,
		storedUser.ID,
	)
	if err != nil {
		return nil, false
	}
	return storedUser, true
}

func (r *UserRepository) GetPreferences(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT preference_value
FROM user_preferences
WHERE user_id = ?
ORDER BY id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	return items, rows.Err()
}

func (r *UserRepository) SavePreferences(ctx context.Context, userID string, preferences []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_preferences WHERE user_id = ?`, userID); err != nil {
		return err
	}
	for _, item := range preferences {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO user_preferences(user_id, preference_value)
VALUES(?, ?)`, userID, item); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *UserRepository) GetSettings(ctx context.Context, userID string) (*repository.UserSettings, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT trip_reminder_enabled, price_reminder_enabled, personalized_recommend_enabled
FROM user_settings
WHERE user_id = ?
LIMIT 1`, userID)
	var settings repository.UserSettings
	var tripReminder int
	var priceReminder int
	var personalized int
	err := row.Scan(&tripReminder, &priceReminder, &personalized)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	settings.TripReminderEnabled = tripReminder == 1
	settings.PriceReminderEnabled = priceReminder == 1
	settings.PersonalizedRecommendEnabled = personalized == 1
	return &settings, nil
}

func (r *UserRepository) SaveSettings(ctx context.Context, userID string, settings *repository.UserSettings) error {
	if settings == nil {
		return nil
	}
	_, err := r.db.ExecContext(ctx, `
INSERT INTO user_settings(
    user_id, trip_reminder_enabled, price_reminder_enabled, personalized_recommend_enabled
) VALUES(?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    trip_reminder_enabled = VALUES(trip_reminder_enabled),
    price_reminder_enabled = VALUES(price_reminder_enabled),
    personalized_recommend_enabled = VALUES(personalized_recommend_enabled)`,
		userID,
		userBoolToInt(settings.TripReminderEnabled),
		userBoolToInt(settings.PriceReminderEnabled),
		userBoolToInt(settings.PersonalizedRecommendEnabled),
	)
	return err
}

func (r *UserRepository) CountTrips(ctx context.Context, userID string) (int64, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT COUNT(1)
FROM trips
WHERE user_id = ? AND deleted_at IS NULL`, userID)
	var count int64
	if err := row.Scan(&count); err != nil {
		if stringsContainsTableError(err, "trips") {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) findOne(ctx context.Context, where string, arg any) (*repository.User, bool) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, phone, password_hash, nickname, avatar_url, home_city, current_city, status, created_at, updated_at, deleted_at
FROM users
WHERE `+where+` AND deleted_at IS NULL
LIMIT 1`, arg)
	var item repository.User
	var id int64
	var phone sql.NullString
	var deletedAt sql.NullTime
	err := row.Scan(
		&id,
		&phone,
		&item.PasswordHash,
		&item.Nickname,
		&item.AvatarURL,
		&item.HomeCity,
		&item.CurrentCity,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	item.ID = strconv.FormatInt(id, 10)
	item.Phone = phone.String
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}
	return &item, true
}

func cloneUser(input *repository.User) *repository.User {
	if input == nil {
		return nil
	}
	cloned := *input
	return &cloned
}

func userBoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func stringsContainsTableError(err error, tableName string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tableName))
}

func isDuplicateColumnError(err error, columnName string) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate column") && strings.Contains(message, strings.ToLower(columnName))
}
