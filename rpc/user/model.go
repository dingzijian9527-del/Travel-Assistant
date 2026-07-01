package rpcuser

import (
	"regexp"
	"strings"
	"time"
)

type userModel struct {
	ID            string
	Phone         string
	PasswordHash  string
	Nickname      string
	AvatarURL     string
	HomeCity      string
	CurrentCity   string
	MemberLevel   string
	AccountStatus int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
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
