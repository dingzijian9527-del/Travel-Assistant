package handler

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	errPasswordTooShort   = errors.New("密码不能少于8位")
	errPasswordTooLong    = errors.New("密码不能超过64位")
	errPasswordMissingNum = errors.New("密码必须包含数字")
	errPasswordMissingLow = errors.New("密码必须包含小写字母")
	errPasswordMissingUp  = errors.New("密码必须包含大写字母")
)

func validatePassword(password string) error {
	trimmed := strings.TrimSpace(password)
	length := utf8.RuneCountInString(trimmed)
	if length < 8 {
		return errPasswordTooShort
	}
	if length > 64 {
		return errPasswordTooLong
	}

	var hasDigit, hasLower, hasUpper bool
	for _, ch := range trimmed {
		switch {
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		}
	}

	if !hasDigit {
		return errPasswordMissingNum
	}
	if !hasLower {
		return errPasswordMissingLow
	}
	if !hasUpper {
		return errPasswordMissingUp
	}
	return nil
}
