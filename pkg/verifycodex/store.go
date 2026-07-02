package verifycodex

import (
	"context"
	"crypto/subtle"
	"strings"
	"sync"
	"time"
)

// Store 定义验证码保存和校验能力。
type Store interface {
	Save(ctx context.Context, phone string, code string, ttl time.Duration) error
	Check(ctx context.Context, phone string, code string) (bool, error)
	Verify(ctx context.Context, phone string, code string) (bool, error)
	Delete(ctx context.Context, phone string) error
}

type memoryItem struct {
	code     string
	expireAt time.Time
}

// MemoryStore 是单元测试使用的内存验证码存储。
type MemoryStore struct {
	mu    sync.Mutex
	items map[string]memoryItem
}

// NewMemoryStore 创建内存验证码存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]memoryItem{}}
}

// Save 保存验证码和过期时间。
func (s *MemoryStore) Save(_ context.Context, phone string, code string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key(phone)] = memoryItem{code: strings.TrimSpace(code), expireAt: time.Now().Add(ttl)}
	return nil
}

// Verify 校验验证码，校验成功后立即删除，避免重复使用。
func (s *MemoryStore) Verify(_ context.Context, phone string, code string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	matched := s.checkLocked(phone, code)
	if !matched {
		return false, nil
	}
	delete(s.items, key(phone))
	return true, nil
}

// Check 只校验验证码是否匹配，不删除验证码。
func (s *MemoryStore) Check(_ context.Context, phone string, code string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.checkLocked(phone, code), nil
}

func (s *MemoryStore) checkLocked(phone string, code string) bool {
	cacheKey := key(phone)
	item, ok := s.items[cacheKey]
	if !ok || time.Now().After(item.expireAt) {
		delete(s.items, cacheKey)
		return false
	}
	if subtle.ConstantTimeCompare([]byte(item.code), []byte(strings.TrimSpace(code))) != 1 {
		return false
	}
	return true
}

// Delete 删除指定手机号的验证码。
func (s *MemoryStore) Delete(_ context.Context, phone string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key(phone))
	return nil
}

func key(phone string) string {
	return "travel-assistant:user:register-code:" + strings.TrimSpace(phone)
}
