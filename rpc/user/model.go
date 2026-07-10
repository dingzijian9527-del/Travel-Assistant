package rpcuser

import (
	"regexp"
	"strings"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type userModel = repository.User
type userSettingsModel = repository.UserSettings

type userDashboardModel struct {
	User        *userModel
	Stats       *userStatsModel
	Preferences []string
	Settings    *userSettingsModel
}

type userStatsModel struct {
	TripCount     int64
	FavoriteCount int64
	UnreadCount   int64
	CouponCount   int64
}

func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

func validPhone(phone string) bool {
	return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(phone)
}

func defaultNickname(phone string, nickname *string) string {
	if nickname != nil && strings.TrimSpace(*nickname) != "" {
		return strings.TrimSpace(*nickname)
	}
	if len(phone) >= 4 {
		return "用户" + phone[len(phone)-4:]
	}
	return "旅行者"
}

func optionalTrimmed(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func defaultUserSettings() *userSettingsModel {
	return &userSettingsModel{
		TripReminderEnabled:          true,
		PriceReminderEnabled:         true,
		PersonalizedRecommendEnabled: false,
	}
}

func normalizePreferences(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
