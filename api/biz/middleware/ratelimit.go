package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type windowCounter struct {
	count   int
	resetAt time.Time
}

type FixedWindowLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]windowCounter
	now     func() time.Time
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &FixedWindowLimiter{
		limit:   limit,
		window:  window,
		entries: make(map[string]windowCounter),
		now:     time.Now,
	}
}

func (l *FixedWindowLimiter) Allow(key string) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return true
	}
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, ok := l.entries[key]
	if !ok || !now.Before(entry.resetAt) {
		l.entries[key] = windowCounter{count: 1, resetAt: now.Add(l.window)}
		return true
	}
	if entry.count >= l.limit {
		return false
	}
	entry.count++
	l.entries[key] = entry
	return true
}

func RateLimitByIP(name string, limit int, window time.Duration) app.HandlerFunc {
	limiter := NewFixedWindowLimiter(limit, window)
	return rateLimit(name, limiter, func(c *app.RequestContext) string {
		return "ip:" + clientIP(c)
	}, "请求过于频繁，请稍后再试")
}

func RegisterCodeRateLimit() app.HandlerFunc {
	phoneLimiter := NewFixedWindowLimiter(1, time.Minute)
	deviceLimiter := NewFixedWindowLimiter(3, 10*time.Minute)
	ipLimiter := NewFixedWindowLimiter(5, 10*time.Minute)
	return func(ctx context.Context, c *app.RequestContext) {
		payload := struct {
			Phone string `json:"phone"`
		}{}
		body := c.Request.Body()
		if len(body) > 0 {
			_ = json.Unmarshal(body, &payload)
		}
		phone := normalizePhone(payload.Phone)
		if phone != "" && !phoneLimiter.Allow("phone:"+phone) {
			abortTooManyRequests(c, "验证码发送过于频繁，请稍后再试")
			return
		}
		deviceID := strings.TrimSpace(string(c.GetHeader("X-Device-ID")))
		if deviceID != "" && !deviceLimiter.Allow("device:"+deviceID) {
			abortTooManyRequests(c, "验证码发送过于频繁，请稍后再试")
			return
		}
		if !ipLimiter.Allow("ip:" + clientIP(c)) {
			abortTooManyRequests(c, "验证码发送过于频繁，请稍后再试")
			return
		}
		c.Next(ctx)
	}
}

func rateLimit(name string, limiter *FixedWindowLimiter, keyFn func(*app.RequestContext) string, message string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if limiter == nil || keyFn == nil {
			c.Next(ctx)
			return
		}
		key := strings.TrimSpace(keyFn(c))
		if key == "" || limiter.Allow(name+":"+key) {
			c.Next(ctx)
			return
		}
		abortTooManyRequests(c, message)
	}
}

func abortTooManyRequests(c *app.RequestContext, message string) {
	c.AbortWithStatusJSON(consts.StatusTooManyRequests, map[string]any{
		"code": consts.StatusTooManyRequests,
		"msg":  message,
	})
}

func clientIP(c *app.RequestContext) string {
	forwarded := strings.TrimSpace(string(c.GetHeader("X-Forwarded-For")))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}
	if ip := strings.TrimSpace(c.ClientIP()); ip != "" {
		return ip
	}
	return "unknown"
}

func normalizePhone(phone string) string {
	replacer := strings.NewReplacer(" ", "", "-", "")
	return strings.TrimSpace(replacer.Replace(phone))
}
