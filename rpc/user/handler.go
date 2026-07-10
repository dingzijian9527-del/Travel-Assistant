package rpcuser

import (
	"context"
	"sync"

	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
)

// UserServiceImpl 实现接口定义中的用户服务。
type UserServiceImpl struct {
	initOnce sync.Once
	service  *userService
	initErr  error
}

func (s *UserServiceImpl) getService() (*userService, error) {
	s.initOnce.Do(func() {
		cfg := config.MustGlobal()
		repo, err := mysqlrepo.NewUserRepository(cfg.MySQL)
		if err != nil {
			s.initErr = err
			return
		}
		s.service = newUserServiceWithAuth(repo, cfg.Auth)
	})
	return s.service, s.initErr
}

// Register 实现用户注册接口。
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	if req == nil {
		return registerErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return registerErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	createdUser, svcErr := service.Register(ctx, req)
	if svcErr != nil {
		return registerErrorResponse(svcErr), nil
	}
	return &user.RegisterResponse{BaseResp: successBaseResp(), User: toUserDTO(createdUser)}, nil
}

// Login 实现用户登录接口。
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	if req == nil {
		return loginErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return loginErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	loggedInUser, token, svcErr := service.Login(ctx, req)
	if svcErr != nil {
		return loginErrorResponse(svcErr), nil
	}
	return &user.LoginResponse{BaseResp: successBaseResp(), Token: &token, User: toUserDTO(loggedInUser)}, nil
}

// GetProfile 实现用户资料查询接口。
func (s *UserServiceImpl) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	if req == nil {
		return profileErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return profileErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	profile, svcErr := service.GetProfile(ctx, req.Id)
	if svcErr != nil {
		return profileErrorResponse(svcErr), nil
	}
	return &user.GetProfileResponse{BaseResp: successBaseResp(), User: toUserDTO(profile)}, nil
}

// UpdateProfile 实现用户资料更新接口。
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, req *user.UpdateProfileRequest) (*user.UpdateProfileResponse, error) {
	if req == nil {
		return updateProfileErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return updateProfileErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	updatedUser, svcErr := service.UpdateProfile(ctx, req)
	if svcErr != nil {
		return updateProfileErrorResponse(svcErr), nil
	}
	return &user.UpdateProfileResponse{BaseResp: successBaseResp(), User: toUserDTO(updatedUser)}, nil
}

func (s *UserServiceImpl) GetDashboard(ctx context.Context, req *user.GetDashboardRequest) (*user.GetDashboardResponse, error) {
	if req == nil {
		return dashboardErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return dashboardErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	dashboard, svcErr := service.GetDashboard(ctx, req.Id)
	if svcErr != nil {
		return dashboardErrorResponse(svcErr), nil
	}
	return &user.GetDashboardResponse{BaseResp: successBaseResp(), Dashboard: toUserDashboardDTO(dashboard)}, nil
}

func (s *UserServiceImpl) GetPreferences(ctx context.Context, req *user.GetPreferencesRequest) (*user.GetPreferencesResponse, error) {
	if req == nil {
		return preferencesErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return preferencesErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	items, svcErr := service.GetPreferences(ctx, req.Id)
	if svcErr != nil {
		return preferencesErrorResponse(svcErr), nil
	}
	return &user.GetPreferencesResponse{BaseResp: successBaseResp(), Preferences: toUserPreferencesDTO(items)}, nil
}

func (s *UserServiceImpl) UpdatePreferences(ctx context.Context, req *user.UpdatePreferencesRequest) (*user.UpdatePreferencesResponse, error) {
	if req == nil {
		return updatePreferencesErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return updatePreferencesErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	items, svcErr := service.UpdatePreferences(ctx, req.Id, req.Items)
	if svcErr != nil {
		return updatePreferencesErrorResponse(svcErr), nil
	}
	return &user.UpdatePreferencesResponse{BaseResp: successBaseResp(), Preferences: toUserPreferencesDTO(items)}, nil
}

func (s *UserServiceImpl) GetSettings(ctx context.Context, req *user.GetSettingsRequest) (*user.GetSettingsResponse, error) {
	if req == nil {
		return settingsErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return settingsErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	settings, svcErr := service.GetSettings(ctx, req.Id)
	if svcErr != nil {
		return settingsErrorResponse(svcErr), nil
	}
	return &user.GetSettingsResponse{BaseResp: successBaseResp(), Settings: toUserSettingsDTO(settings)}, nil
}

func (s *UserServiceImpl) UpdateSettings(ctx context.Context, req *user.UpdateSettingsRequest) (*user.UpdateSettingsResponse, error) {
	if req == nil {
		return updateSettingsErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService()
	if err != nil {
		return updateSettingsErrorResponse(errInternal("用户服务初始化失败")), nil
	}
	settings, svcErr := service.UpdateSettings(ctx, req.Id, &userSettingsModel{
		TripReminderEnabled:          req.TripReminderEnabled,
		PriceReminderEnabled:         req.PriceReminderEnabled,
		PersonalizedRecommendEnabled: req.PersonalizedRecommendEnabled,
	})
	if svcErr != nil {
		return updateSettingsErrorResponse(svcErr), nil
	}
	return &user.UpdateSettingsResponse{BaseResp: successBaseResp(), Settings: toUserSettingsDTO(settings)}, nil
}
