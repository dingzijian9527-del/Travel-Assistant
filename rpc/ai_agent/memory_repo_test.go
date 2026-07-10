package rpcaiagent

import (
	"context"
	"sync"
)

type aiAgentRepo struct {
	mu              sync.RWMutex
	messagesByUser  map[int64][]chatMessageModel
	maxHistoryItems int
}

func newAIAgentRepo() *aiAgentRepo {
	return &aiAgentRepo{messagesByUser: make(map[int64][]chatMessageModel), maxHistoryItems: 20}
}

func (r *aiAgentRepo) AppendMessages(_ context.Context, userID int64, messages ...*chatMessageModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, message := range messages {
		if message == nil {
			continue
		}
		r.messagesByUser[userID] = append(r.messagesByUser[userID], *message)
	}
	if len(r.messagesByUser[userID]) > r.maxHistoryItems {
		r.messagesByUser[userID] = r.messagesByUser[userID][len(r.messagesByUser[userID])-r.maxHistoryItems:]
	}
}

func (r *aiAgentRepo) History(_ context.Context, userID int64) []chatMessageModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.messagesByUser[userID]
	copied := make([]chatMessageModel, len(items))
	copy(copied, items)
	return copied
}
