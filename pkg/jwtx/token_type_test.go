package jwtx

import (
	"testing"
	"time"
)

func TestRefreshTokenCannotBeUsedAsAccessToken(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	refresh, err := GenerateRefresh(cfg, Claims{UserID: "user-1"}, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("生成刷新令牌失败：%v", err)
	}
	if _, err := Parse(cfg, refresh); err == nil {
		t.Fatal("刷新令牌不应通过普通访问令牌解析")
	}
	if _, err := ParseRefresh(cfg, refresh); err != nil {
		t.Fatalf("刷新令牌应通过专用解析：%v", err)
	}
}

func TestAccessTokenCannotBeUsedAsRefreshToken(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	access, err := Generate(cfg, Claims{UserID: "user-2"})
	if err != nil {
		t.Fatalf("生成访问令牌失败：%v", err)
	}
	if _, err := ParseRefresh(cfg, access); err == nil {
		t.Fatal("访问令牌不应通过刷新令牌解析")
	}
}
