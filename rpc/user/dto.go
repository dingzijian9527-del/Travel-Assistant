package rpcuser

import (
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"
	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
)

func toUserDTO(input *userModel) *user.UserInfo {
	if input == nil {
		return nil
	}
	return &user.UserInfo{
		Id:            input.ID,
		Phone:         input.Phone,
		Nickname:      input.Nickname,
		AvatarUrl:     optionalString(input.AvatarURL),
		HomeCity:      optionalString(input.HomeCity),
		CurrentCity:   optionalString(input.CurrentCity),
		MemberLevel:   "",
		AccountStatus: input.Status,
		CreatedAt:     optionalTime(input.CreatedAt),
		UpdatedAt:     optionalTime(input.UpdatedAt),
		DeletedAt:     optionalTimePtr(input.DeletedAt),
	}
}

func toUserSettingsDTO(input *userSettingsModel) *user.UserSettings {
	if input == nil {
		return nil
	}
	return &user.UserSettings{
		TripReminderEnabled:          input.TripReminderEnabled,
		PriceReminderEnabled:         input.PriceReminderEnabled,
		PersonalizedRecommendEnabled: input.PersonalizedRecommendEnabled,
	}
}

func toUserPreferencesDTO(items []string) *user.UserPreferences {
	return &user.UserPreferences{Items: normalizePreferences(items)}
}

func toUserStatsDTO(input *userStatsModel) *user.UserStats {
	if input == nil {
		return nil
	}
	return &user.UserStats{
		TripCount:     input.TripCount,
		FavoriteCount: input.FavoriteCount,
		UnreadCount:   input.UnreadCount,
		CouponCount:   input.CouponCount,
	}
}

func toUserDashboardDTO(input *userDashboardModel) *user.UserDashboard {
	if input == nil {
		return nil
	}
	return &user.UserDashboard{
		User:        toUserDTO(input.User),
		Stats:       toUserStatsDTO(input.Stats),
		Settings:    toUserSettingsDTO(input.Settings),
		Preferences: toUserPreferencesDTO(input.Preferences),
	}
}

func successBaseResp() *base.BaseResp {
	return &base.BaseResp{Code: int32(base.ErrorCode_SUCCESS), Msg: "success"}
}

func errorBaseResp(err *serviceError) *base.BaseResp {
	if err == nil {
		return successBaseResp()
	}
	return &base.BaseResp{Code: int32(err.code), Msg: err.message}
}

func registerErrorResponse(err *serviceError) *user.RegisterResponse {
	return &user.RegisterResponse{BaseResp: errorBaseResp(err)}
}

func loginErrorResponse(err *serviceError) *user.LoginResponse {
	return &user.LoginResponse{BaseResp: errorBaseResp(err)}
}

func profileErrorResponse(err *serviceError) *user.GetProfileResponse {
	return &user.GetProfileResponse{BaseResp: errorBaseResp(err)}
}

func updateProfileErrorResponse(err *serviceError) *user.UpdateProfileResponse {
	return &user.UpdateProfileResponse{BaseResp: errorBaseResp(err)}
}

func dashboardErrorResponse(err *serviceError) *user.GetDashboardResponse {
	return &user.GetDashboardResponse{BaseResp: errorBaseResp(err)}
}

func preferencesErrorResponse(err *serviceError) *user.GetPreferencesResponse {
	return &user.GetPreferencesResponse{BaseResp: errorBaseResp(err)}
}

func updatePreferencesErrorResponse(err *serviceError) *user.UpdatePreferencesResponse {
	return &user.UpdatePreferencesResponse{BaseResp: errorBaseResp(err)}
}

func settingsErrorResponse(err *serviceError) *user.GetSettingsResponse {
	return &user.GetSettingsResponse{BaseResp: errorBaseResp(err)}
}

func updateSettingsErrorResponse(err *serviceError) *user.UpdateSettingsResponse {
	return &user.UpdateSettingsResponse{BaseResp: errorBaseResp(err)}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func optionalTime(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.Format("2006-01-02 15:04:05")
	return &formatted
}

func optionalTimePtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	return optionalTime(*value)
}
