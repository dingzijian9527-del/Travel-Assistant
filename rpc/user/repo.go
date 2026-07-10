package rpcuser

import "context"

type userRepository interface {
	Create(ctx context.Context, input *userModel) (*userModel, error)
	FindByPhone(ctx context.Context, phone string) (*userModel, bool)
	FindByID(ctx context.Context, userID string) (*userModel, bool)
	Update(ctx context.Context, userID string, apply func(*userModel)) (*userModel, bool)
	GetPreferences(ctx context.Context, userID string) ([]string, error)
	SavePreferences(ctx context.Context, userID string, preferences []string) error
	GetSettings(ctx context.Context, userID string) (*userSettingsModel, error)
	SaveSettings(ctx context.Context, userID string, settings *userSettingsModel) error
	CountTrips(ctx context.Context, userID string) (int64, error)
}
