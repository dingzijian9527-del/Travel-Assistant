package rpcuser

import (
	"context"
	"sync"

	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
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
		repo, err := newMySQLUserRepo(cfg.MySQL)
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
