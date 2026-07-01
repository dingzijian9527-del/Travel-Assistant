package validator

import "strings"

// RequiredString 校验必填字符串。
func RequiredString(value string, field string) (string, bool, string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false, field + "不能为空"
	}
	return trimmed, true, ""
}
