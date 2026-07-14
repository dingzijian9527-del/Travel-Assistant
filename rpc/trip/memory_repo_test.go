package rpctrip

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"
)

type tripRepo struct {
	mu     sync.RWMutex
	nextID int64
	trips  map[string]*tripModel
}

func newTripRepo() *tripRepo {
	return &tripRepo{trips: make(map[string]*tripModel)}
}

func (r *tripRepo) Create(_ context.Context, input *tripModel) (*tripModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	stored := cloneTrip(input)
	stored.ID = strconv.FormatInt(r.nextID, 10)
	r.trips[stored.ID] = cloneTrip(stored)
	return cloneTrip(stored), nil
}

func (r *tripRepo) ListByUser(_ context.Context, userID int64, limit int) ([]*tripModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*tripModel, 0)
	for _, item := range r.trips {
		if item.UserID == userID && item.DeletedAt == nil {
			items = append(items, cloneTrip(item))
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (r *tripRepo) LatestByUser(ctx context.Context, userID int64) (*tripModel, bool, error) {
	items, err := r.ListByUser(ctx, userID, 1)
	if err != nil {
		return nil, false, err
	}
	if len(items) == 0 {
		return nil, false, nil
	}
	return items[0], true, nil
}

func (r *tripRepo) FindByIDForUser(_ context.Context, userID int64, tripID string) (*tripModel, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.trips[tripID]
	if !ok || item.UserID != userID || item.DeletedAt != nil {
		return nil, false, nil
	}
	return cloneTrip(item), true, nil
}

func (r *tripRepo) DeleteByIDForUser(_ context.Context, userID int64, tripID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.trips[tripID]
	if !ok || item.UserID != userID || item.DeletedAt != nil {
		return false
	}
	now := time.Now()
	item.DeletedAt = &now
	item.UpdatedAt = now
	return true
}

func (r *tripRepo) Update(_ context.Context, userID int64, tripID string, apply func(*tripModel)) (*tripModel, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.trips[tripID]
	if !ok || item.UserID != userID || item.DeletedAt != nil {
		return nil, false, nil
	}
	apply(item)
	item.UpdatedAt = time.Now()
	r.trips[tripID] = item
	return cloneTrip(item), true, nil
}

func cloneTrip(input *tripModel) *tripModel {
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
