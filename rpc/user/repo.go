package rpcuser

import (
	"context"
	"strconv"
	"sync"
)

type userRepo struct {
	mu        sync.RWMutex
	nextID    int64
	usersByID map[string]*userModel
	phoneToID map[string]string
}

type userRepository interface {
	Create(ctx context.Context, input *userModel) (*userModel, error)
	FindByPhone(ctx context.Context, phone string) (*userModel, bool)
	FindByID(ctx context.Context, userID string) (*userModel, bool)
	Update(ctx context.Context, userID string, apply func(*userModel)) (*userModel, bool)
}

func newUserRepo() *userRepo {
	return &userRepo{
		nextID:    1,
		usersByID: make(map[string]*userModel),
		phoneToID: make(map[string]string),
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

func newUserID(sequence int64) string {
	return "user-" + strconv.FormatInt(sequence, 10)
}

func cloneUser(input *userModel) *userModel {
	if input == nil {
		return nil
	}
	cloned := *input
	return &cloned
}
