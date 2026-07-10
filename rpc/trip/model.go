package rpctrip

import (
	"strings"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type tripModel = repository.Trip

func defaultTripTitle(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return "我的旅行行程"
	}
	return trimmed
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func optionalInt32(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}
