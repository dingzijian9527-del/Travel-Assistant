package redisx

import (
	"errors"
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestNewRejectsEmptyAddr(t *testing.T) {
	_, err := New(config.RedisConfig{})
	if !errors.Is(err, ErrMissingAddr) {
		t.Fatalf("空缓存地址应返回 ErrMissingAddr，实际: %v", err)
	}
}

func TestOptionsFromConfig(t *testing.T) {
	options := Options(config.RedisConfig{
		Addr:     "127.0.0.1:6379",
		Username: "tester",
		Password: "secret",
		DB:       2,
	})

	if options.Addr != "127.0.0.1:6379" || options.Username != "tester" || options.Password != "secret" || options.DB != 2 {
		t.Fatalf("缓存连接参数转换异常: %#v", options)
	}
}
