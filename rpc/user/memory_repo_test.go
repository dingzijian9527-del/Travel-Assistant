package rpcuser

import (
	"context"
	"strconv"
	"sync"
)

type userRepo struct {
	mu              sync.RWMutex
	nextID          int64
	usersByID       map[string]*userModel
	phoneToID       map[string]string
	preferencesByID map[string][]string
	settingsByID    map[string]*userSettingsModel
}

func newUserRepo() *userRepo {
	return &userRepo{
		nextID:    1,
		usersByID: make(map[string]*userModel),
		phoneToID: make(map[string]string),
		preferencesByID: make(map[string][]string),
		settingsByID: make(map[string]*userSettingsModel),
	}
}

func (r *userRepo) Create(_ context.Context, input *userModel) (*userModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	storedUser := cloneUser(input)
	storedUser.ID = newUserID(r.nextID)
	r.nextID++
	r.usersByID[storedUser.ID] = cloneUser(storedUser)
	r.phoneToID[storedUser.Phone] = storedUser.ID
	return cloneUser(storedUser), nil
}

func (r *userRepo) FindByPhone(_ context.Context, phone string) (*userModel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userID, exists := r.phoneToID[phone]
	if !exists {
		return nil, false
	}
	return cloneUser(r.usersByID[userID]), true
}

func (r *userRepo) FindByID(_ context.Context, userID string) (*userModel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storedUser, exists := r.usersByID[userID]
	if !exists {
		return nil, false
	}
	return cloneUser(storedUser), true
}

func (r *userRepo) Update(_ context.Context, userID string, apply func(*userModel)) (*userModel, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	storedUser, exists := r.usersByID[userID]
	if !exists {
		return nil, false
	}
	apply(storedUser)
	return cloneUser(storedUser), true
}

func (r *userRepo) GetPreferences(_ context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.preferencesByID[userID]
	return append([]string(nil), items...), nil
}

func (r *userRepo) SavePreferences(_ context.Context, userID string, preferences []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.preferencesByID[userID] = append([]string(nil), preferences...)
	return nil
}

func (r *userRepo) GetSettings(_ context.Context, userID string) (*userSettingsModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored := r.settingsByID[userID]
	if stored == nil {
		return nil, nil
	}
	cloned := *stored
	return &cloned, nil
}

func (r *userRepo) SaveSettings(_ context.Context, userID string, settings *userSettingsModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if settings == nil {
		delete(r.settingsByID, userID)
		return nil
	}
	cloned := *settings
	r.settingsByID[userID] = &cloned
	return nil
}

func (r *userRepo) CountTrips(_ context.Context, _ string) (int64, error) {
	return 0, nil
}

func newUserID(sequence int64) string {
	return strconv.FormatInt(sequence, 10)
}

func cloneUser(input *userModel) *userModel {
	if input == nil {
		return nil
	}
	cloned := *input
	return &cloned
}
