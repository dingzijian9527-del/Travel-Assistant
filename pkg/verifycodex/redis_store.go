package verifycodex

import (
	"context"
	"crypto/subtle"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore 使用缓存服务保存注册验证码。
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore 创建缓存服务验证码存储。
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

// Save 将验证码写入缓存服务，并设置过期时间。
func (s *RedisStore) Save(ctx context.Context, phone string, code string, ttl time.Duration) error {
	return s.client.Set(ctx, key(phone), strings.TrimSpace(code), ttl).Err()
}

// Verify 从缓存服务读取验证码并校验，成功后删除。
func (s *RedisStore) Verify(ctx context.Context, phone string, code string) (bool, error) {
	ok, err := s.Check(ctx, phone, code)
	if err != nil || !ok {
		return ok, err
	}
	if err := s.client.Del(ctx, key(phone)).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// Check 从缓存服务读取验证码并校验，不删除验证码。
func (s *RedisStore) Check(ctx context.Context, phone string, code string) (bool, error) {
	cacheKey := key(phone)
	expected, err := s.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(expected), []byte(strings.TrimSpace(code))) != 1 {
		return false, nil
	}
	return true, nil
}

// Delete 删除指定手机号的验证码。
func (s *RedisStore) Delete(ctx context.Context, phone string) error {
	return s.client.Del(ctx, key(phone)).Err()
}
