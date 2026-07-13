package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	maxLoginFailures   = 5
	loginFailureWindow = 15 * time.Minute
)

type loginAttempt struct {
	failures  int
	firstFail time.Time
}

// LoginRateLimit 限制同一手机号的登录失败频率，超出后返回统一错误。
// 生产环境应改用 Redis 实现跨进程共享。
func legacyLoginRateLimit() app.HandlerFunc {
	mu := &sync.Mutex{}
	attempts := map[string]*loginAttempt{}

	go func() {
		for {
			time.Sleep(loginFailureWindow)
			mu.Lock()
			now := time.Now()
			for key, attempt := range attempts {
				if now.Sub(attempt.firstFail) > loginFailureWindow {
					delete(attempts, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(ctx context.Context, c *app.RequestContext) {
		phone := loginPhoneFromRequest(c)
		if phone == "" {
			c.Next(ctx)
			return
		}

		fingerprint := hashLoginFingerprint(phone, c)
		mu.Lock()
		entry, exists := attempts[fingerprint]
		if exists && entry.failures >= maxLoginFailures {
			if time.Since(entry.firstFail) < loginFailureWindow {
				mu.Unlock()
				c.AbortWithStatusJSON(consts.StatusTooManyRequests, map[string]any{
					"code": 429,
					"msg":  "登录次数过多，请稍后再试",
				})
				return
			}
			delete(attempts, fingerprint)
		}
		mu.Unlock()

		c.Next(ctx)

		// 登录失败后计数
		if c.Response.StatusCode() != consts.StatusOK {
			mu.Lock()
			entry, exists = attempts[fingerprint]
			if !exists || time.Since(entry.firstFail) > loginFailureWindow {
				attempts[fingerprint] = &loginAttempt{failures: 1, firstFail: time.Now()}
			} else {
				entry.failures++
			}
			mu.Unlock()
		}
	}
}

func loginPhoneFromRequest(c *app.RequestContext) string {
	// 仅对登录接口提取手机号用于限流
	if string(c.Request.Path()) != "/api/v1/user/login" {
		return ""
	}
	body := c.Request.Body()
	if len(body) == 0 {
		return ""
	}
	// 简单提取 JSON 中的 phone 字段
	raw := string(body)
	idx := strings.Index(raw, `"phone"`)
	if idx < 0 {
		return ""
	}
	idx = strings.Index(raw[idx:], `"`)
	if idx < 0 {
		return ""
	}
	idx += len(`"`)
	rest := raw[strings.Index(raw, `"phone"`)+len(`"phone"`):]
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
	clientIP := string(c.Request.Header.Peek("X-Forwarded-For"))
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	source := phone + "|" + clientIP
	hash := sha256.Sum256([]byte(source))
	return hex.EncodeToString(hash[:])
}
