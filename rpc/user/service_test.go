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
	if createdUser.HomeCity != "成都" || createdUser.CurrentCity != "上海" || createdUser.Status != 1 {
		t.Fatalf("created user fields do not match current users table fields: %+v", createdUser)
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
	if profile.Nickname != nickname || profile.CurrentCity != "上海" {
		t.Fatalf("unexpected profile: %+v", profile)
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
		t.Fatalf("updated profile fields do not match current users table fields: %+v", updatedProfile)
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

func TestUserServiceDashboardPreferencesAndSettings(t *testing.T) {
	service := newUserService(newUserRepo())
	createdUser, svcErr := service.Register(context.Background(), &user.RegisterRequest{
		Phone:       "13800138002",
		Password:    "safe-pass-123",
		Nickname:    stringPtr("测试用户"),
		HomeCity:    stringPtr("成都"),
		CurrentCity: stringPtr("杭州"),
	})
	if svcErr != nil {
		t.Fatalf("register returned error: %v", svcErr)
	}

	preferences, svcErr := service.GetPreferences(context.Background(), createdUser.ID)
	if svcErr != nil {
		t.Fatalf("get preferences returned error: %v", svcErr)
	}
	if len(preferences) != 0 {
		t.Fatalf("expected empty preferences by default, got: %#v", preferences)
	}

	savedPreferences, svcErr := service.UpdatePreferences(context.Background(), createdUser.ID, []string{"川菜", "茶馆", "川菜", "  海岛  "})
	if svcErr != nil {
		t.Fatalf("update preferences returned error: %v", svcErr)
	}
	if len(savedPreferences) != 3 || savedPreferences[0] != "川菜" || savedPreferences[2] != "海岛" {
		t.Fatalf("unexpected saved preferences: %#v", savedPreferences)
	}

	settings, svcErr := service.GetSettings(context.Background(), createdUser.ID)
	if svcErr != nil {
		t.Fatalf("get settings returned error: %v", svcErr)
	}
	if !settings.TripReminderEnabled || !settings.PriceReminderEnabled || settings.PersonalizedRecommendEnabled {
		t.Fatalf("unexpected default settings: %#v", settings)
	}

	updatedSettings, svcErr := service.UpdateSettings(context.Background(), createdUser.ID, &userSettingsModel{
		TripReminderEnabled:          false,
		PriceReminderEnabled:         true,
		PersonalizedRecommendEnabled: true,
	})
	if svcErr != nil {
		t.Fatalf("update settings returned error: %v", svcErr)
	}
	if updatedSettings.TripReminderEnabled || !updatedSettings.PriceReminderEnabled || !updatedSettings.PersonalizedRecommendEnabled {
		t.Fatalf("unexpected updated settings: %#v", updatedSettings)
	}

	dashboard, svcErr := service.GetDashboard(context.Background(), createdUser.ID)
	if svcErr != nil {
		t.Fatalf("get dashboard returned error: %v", svcErr)
	}
	if dashboard.User == nil || dashboard.User.Nickname != "测试用户" || dashboard.User.CurrentCity != "杭州" {
		t.Fatalf("unexpected dashboard user: %#v", dashboard.User)
	}
	if len(dashboard.Preferences) != 3 {
		t.Fatalf("unexpected dashboard preferences: %#v", dashboard.Preferences)
	}
	if dashboard.Settings == nil || !dashboard.Settings.PersonalizedRecommendEnabled {
		t.Fatalf("unexpected dashboard settings: %#v", dashboard.Settings)
	}
	if dashboard.Stats == nil || dashboard.Stats.TripCount != 0 {
		t.Fatalf("unexpected dashboard stats: %#v", dashboard.Stats)
	}
}

func stringPtr(value string) *string {
	return &value
}
