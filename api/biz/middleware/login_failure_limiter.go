package middleware

import (
	"strings"
	"sync"
	"time"
)

type loginIdentity struct {
	Phone    string
	IP       string
	DeviceID string
}

type LoginFailureLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	attempts map[string]loginAttempt
}

func NewLoginFailureLimiter(limit int, window time.Duration) *LoginFailureLimiter {
	if limit <= 0 {
		limit = 5
	}
	if window <= 0 {
		window = 15 * time.Minute
	}
	return &LoginFailureLimiter{
		limit:    limit,
		window:   window,
		attempts: make(map[string]loginAttempt),
	}
}

func (l *LoginFailureLimiter) Allow(identity loginIdentity) bool {
	keys := loginIdentityKeys(identity)
	if len(keys) == 0 {
		return true
	}

	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, key := range keys {
		entry, ok := l.attempts[key]
		if !ok {
			continue
		}
		if now.Sub(entry.firstFail) >= l.window {
			delete(l.attempts, key)
			continue
		}
		if entry.failures >= l.limit {
			return false
		}
	}
	return true
}

func (l *LoginFailureLimiter) RecordFailure(identity loginIdentity) {
	keys := loginIdentityKeys(identity)
	if len(keys) == 0 {
		return
	}

	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, key := range keys {
		entry, ok := l.attempts[key]
		if !ok || now.Sub(entry.firstFail) >= l.window {
			l.attempts[key] = loginAttempt{failures: 1, firstFail: now}
			continue
		}
		entry.failures++
		l.attempts[key] = entry
	}
}

func (l *LoginFailureLimiter) Reset(identity loginIdentity) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, key := range loginIdentityKeys(identity) {
		delete(l.attempts, key)
	}
}

func loginIdentityKeys(identity loginIdentity) []string {
	keys := make([]string, 0, 3)
	if phone := strings.TrimSpace(identity.Phone); phone != "" {
		keys = append(keys, "phone:"+phone)
	}
	if ip := strings.TrimSpace(identity.IP); ip != "" {
		keys = append(keys, "ip:"+ip)
	}
	if deviceID := strings.TrimSpace(identity.DeviceID); deviceID != "" {
		keys = append(keys, "device:"+deviceID)
	}
	return keys
}
