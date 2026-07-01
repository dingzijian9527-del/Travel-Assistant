package rpcuser

import (
	"context"
	"testing"
	"time"

	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

func TestUserServiceRegisterLoginAndProfile(t *testing.T) {
	service := newUserService(newUserRepo())
	nickname := "旅行者小王"
	createdUser, svcErr := service.Register(context.Background(), &user.RegisterRequest{
		Phone:       "13800138000",
		Password:    "safe-pass-123",
		Nickname:    &nickname,
		HomeCity:    stringPtr("成都"),
		CurrentCity: stringPtr("上海"),
	})
	if svcErr != nil {
		t.Fatalf("register returned error: %v", svcErr)
	}
	if createdUser.ID == "" || createdUser.Phone != "13800138000" {
		t.Fatalf("unexpected created user: %+v", createdUser)
	}
	if createdUser.HomeCity != "成都" || createdUser.CurrentCity != "上海" || createdUser.MemberLevel != "normal" || createdUser.AccountStatus != 1 {
		t.Fatalf("created user fields do not match users table defaults: %+v", createdUser)
	}

	loggedInUser, token, svcErr := service.Login(context.Background(), &user.LoginRequest{
		Phone:    "13800138000",
		Password: "safe-pass-123",
	})
	if svcErr != nil {
		t.Fatalf("login returned error: %v", svcErr)
	}
	if loggedInUser.ID != createdUser.ID || token == "" {
		t.Fatalf("unexpected login result user=%+v token=%q", loggedInUser, token)
	}
	claims, err := jwtx.Parse(jwtx.Config{Secret: "test-jwt-secret", Expire: time.Hour}, token)
	if err != nil {
		t.Fatalf("login token should be a verifiable jwt: %v", err)
	}
	if claims.UserID != createdUser.ID || claims.Phone != createdUser.Phone {
		t.Fatalf("unexpected login token claims: %+v", claims)
	}

	profile, svcErr := service.GetProfile(context.Background(), createdUser.ID)
	if svcErr != nil {
		t.Fatalf("get profile returned error: %v", svcErr)
	}
	if profile.Nickname != nickname {
		t.Fatalf("unexpected profile nickname: %s", profile.Nickname)
	}

	updatedProfile, svcErr := service.UpdateProfile(context.Background(), &user.UpdateProfileRequest{
		Id:          createdUser.ID,
		AvatarUrl:   stringPtr("https://example.test/avatar.png"),
		HomeCity:    stringPtr("重庆"),
		CurrentCity: stringPtr("杭州"),
	})
	if svcErr != nil {
		t.Fatalf("update profile returned error: %v", svcErr)
	}
	if updatedProfile.AvatarURL != "https://example.test/avatar.png" || updatedProfile.HomeCity != "重庆" || updatedProfile.CurrentCity != "杭州" {
		t.Fatalf("updated profile fields do not match users table fields: %+v", updatedProfile)
	}
}

func TestUserServiceRejectsInvalidPassword(t *testing.T) {
	service := newUserService(newUserRepo())
	_, svcErr := service.Register(context.Background(), &user.RegisterRequest{
		Phone:    "13800138001",
		Password: "short",
	})
	if svcErr == nil {
		t.Fatal("expected validation error")
	}
}

func TestUserServiceRejectsInvalidPhone(t *testing.T) {
	service := newUserService(newUserRepo())
	_, svcErr := service.Register(context.Background(), &user.RegisterRequest{
		Phone:    "123",
		Password: "safe-pass-123",
	})
	if svcErr == nil {
		t.Fatal("expected phone validation error")
	}
}

func stringPtr(value string) *string {
	return &value
}
