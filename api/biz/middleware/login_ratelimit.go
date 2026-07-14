package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	maxLoginFailures   = 5
	loginFailureWindow = 15 * time.Minute
)

type loginAttempt struct {
	failures  int
	firstFail time.Time
}

func loginPhoneFromRequest(c *app.RequestContext) string {
	if string(c.Request.Path()) != "/api/v1/user/login" {
		return ""
	}
	body := c.Request.Body()
	if len(body) == 0 {
		return ""
	}
	raw := string(body)
	phoneKey := `"phone"`
	keyIndex := strings.Index(raw, phoneKey)
	if keyIndex < 0 {
		return ""
	}
	rest := raw[keyIndex+len(phoneKey):]
	colon := strings.Index(rest, ":")
	if colon < 0 {
		return ""
	}
	rest = rest[colon+1:]
	start := strings.Index(rest, `"`)
	if start < 0 {
		return ""
	}
	end := strings.Index(rest[start+1:], `"`)
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(rest[start+1 : start+1+end])
}

func hashLoginFingerprint(phone string, c *app.RequestContext) string {
	clientIP := clientIP(c)
	source := phone + "|" + clientIP
	hash := sha256.Sum256([]byte(source))
	return hex.EncodeToString(hash[:])
}
