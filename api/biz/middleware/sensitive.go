package middleware

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const confirmTokenHeader = "X-Confirm-Token"

// ConfirmTokenStore 定义敏感操作确认令牌存储接口。
type ConfirmTokenStore interface {
	Save(ctx context.Context, token string, ttl time.Duration) error
	Verify(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string) error
}

type confirmTokenEntry struct {
	token    string
	expireAt time.Time
}

// memoryConfirmTokenStore 基于内存的确认令牌存储。
type memoryConfirmTokenStore struct {
	mu    sync.Mutex
	items map[string]confirmTokenEntry
}

// NewMemoryConfirmTokenStore 创建内存确认令牌存储。
func NewMemoryConfirmTokenStore() ConfirmTokenStore {
	return &memoryConfirmTokenStore{items: make(map[string]confirmTokenEntry)}
}

var defaultConfirmTokenStore ConfirmTokenStore = NewMemoryConfirmTokenStore()

func DefaultConfirmTokenStore() ConfirmTokenStore {
	return defaultConfirmTokenStore
}

func (s *memoryConfirmTokenStore) Save(_ context.Context, token string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[token] = confirmTokenEntry{token: token, expireAt: time.Now().Add(ttl)}
	return nil
}

func (s *memoryConfirmTokenStore) Verify(_ context.Context, token string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.items[token]
	if !ok || time.Now().After(entry.expireAt) {
		delete(s.items, token)
		return false, nil
	}
	if subtle.ConstantTimeCompare([]byte(entry.token), []byte(token)) != 1 {
		return false, nil
	}
	return true, nil
}

func (s *memoryConfirmTokenStore) Delete(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, token)
	return nil
}

// GenerateConfirmToken 生成一个随机的一次性确认令牌（32字节 hex 字符串）。
func GenerateConfirmToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// RequirePasswordConfirm 要求请求携带有效的 X-Confirm-Token 头。
// 用于保护敏感操作（如修改密码、删除账户等）需要二次校验。
// store 为 nil 时使用内存实现（仅适合单实例部署）。
func RequirePasswordConfirm(store ConfirmTokenStore) app.HandlerFunc {
	if store == nil {
		store = DefaultConfirmTokenStore()
	}

	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.Request.Header.Peek(confirmTokenHeader))
		if token == "" {
			c.AbortWithStatusJSON(consts.StatusForbidden, map[string]any{
				"code": 403,
				"msg":  "需要敏感操作确认，请先调用 /api/v1/user/confirm 获取确认令牌",
			})
			return
		}

		ok, err := store.Verify(ctx, token)
		if err != nil {
			c.AbortWithStatusJSON(consts.StatusInternalServerError, map[string]any{
				"code": 500,
				"msg":  "令牌校验服务异常",
			})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(consts.StatusForbidden, map[string]any{
				"code": 403,
				"msg":  "确认令牌无效或已过期（有效期5分钟），请重新获取",
			})
			return
		}

		// 验证通过后删除令牌（一次性使用）
		_ = store.Delete(ctx, token)
		c.Next(ctx)
	}
}
