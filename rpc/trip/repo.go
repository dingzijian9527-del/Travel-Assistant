package rpctrip

import "context"

type tripRepository interface {
	Create(ctx context.Context, input *tripModel) (*tripModel, error)
	ListByUser(ctx context.Context, userID int64, limit int) ([]*tripModel, error)
	LatestByUser(ctx context.Context, userID int64) (*tripModel, bool, error)
	FindByIDForUser(ctx context.Context, userID int64, tripID string) (*tripModel, bool, error)
	DeleteByIDForUser(ctx context.Context, userID int64, tripID string) bool
}
