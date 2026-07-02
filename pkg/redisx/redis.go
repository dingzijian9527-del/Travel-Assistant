package redisx

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// ErrMissingAddr 表示缺少 Redis 连接地址。
var ErrMissingAddr = errors.New("redis addr 未配置")

// New 创建并探活 Redis 客户端。
func New(cfg config.RedisConfig) (*redis.Client, error) {
	if strings.TrimSpace(cfg.Addr) == "" {
		return nil, ErrMissingAddr
	}
	client := redis.NewClient(Options(cfg))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}

// Options 将项目配置转换为 Redis 客户端参数。
func Options(cfg config.RedisConfig) *redis.Options {
	return &redis.Options{
		Addr:     strings.TrimSpace(cfg.Addr),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}
