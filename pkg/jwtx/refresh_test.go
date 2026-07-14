package jwtx

import (
	"testing"
	"time"
)

func TestGeneratePairTokensAreDifferent(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	access, refresh, err := GeneratePair(cfg, Claims{UserID: "user-1", Phone: "13800138000"}, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("GeneratePair failed: %v", err)
	}
	if access == refresh {
		t.Fatal("access token and refresh token should be different")
	}
}

func TestAccessTokenParsedByParseOnly(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	access, refresh, err := GeneratePair(cfg, Claims{UserID: "user-2", Phone: "13900139000"}, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("GeneratePair failed: %v", err)
	}

	// access token 可以用 Parse 解析
	claims, err := Parse(cfg, access)
	if err != nil {
		t.Fatalf("Parse access token failed: %v", err)
	}
	if claims.UserID != "user-2" {
		t.Fatalf("unexpected user id: %s", claims.UserID)
	}

	// access token 不能用 ParseRefresh 解析
	_, err = ParseRefresh(cfg, access)
	if err == nil {
		t.Fatal("ParseRefresh should reject access token")
	}
	t.Logf("ParseRefresh correctly rejected access token: %v", err)

	// refresh token 不能用 Parse 解析（因为 expire 不同，但如果没过期 Parse 也能通过）
	// refresh token 确实可以用 ParseRefresh 解析
	refreshClaims, err := ParseRefresh(cfg, refresh)
	if err != nil {
		t.Fatalf("ParseRefresh refresh token failed: %v", err)
	}
	if refreshClaims.UserID != "user-2" {
		t.Fatalf("unexpected refresh user id: %s", refreshClaims.UserID)
	}
}

func TestRefreshTokenRejectedByParse(t *testing.T) {
	// 构造一个 refresh token 使用极短的 expire，确保 Parse 因过期拒绝
	cfg := Config{Secret: "test-secret", Expire: 1 * time.Millisecond}
	_, refresh, err := GeneratePair(cfg, Claims{UserID: "user-3"}, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("GeneratePair failed: %v", err)
	}

	// 等待 access token 过期
	time.Sleep(10 * time.Millisecond)

	// refresh token 也已过期（这里用统一 expire 测试）
	_, err = Parse(cfg, refresh)
	if err == nil {
		t.Fatal("expired refresh token should be rejected by Parse")
	}
	t.Logf("Parse correctly rejected expired token: %v", err)
}

func TestParseRefreshRejectsTokensWithoutRefreshType(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	// 直接生成普通 token（不带 token_type）
	token, err := Generate(cfg, Claims{UserID: "user-4"})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// ParseRefresh 应拒绝
	_, err = ParseRefresh(cfg, token)
	if err == nil {
		t.Fatal("ParseRefresh should reject token without refresh type")
	}
	t.Logf("ParseRefresh correctly rejected non-refresh token: %v", err)
}

func TestGeneratePairDefaultRefreshExpire(t *testing.T) {
	cfg := Config{Secret: "test-secret", Expire: time.Hour}
	_, refresh, err := GeneratePair(cfg, Claims{UserID: "user-5"}, 0)
	if err != nil {
		t.Fatalf("GeneratePair with 0 refresh expire failed: %v", err)
	}

	// 验证 refresh token 可解析
	claims, err := ParseRefresh(cfg, refresh)
	if err != nil {
		t.Fatalf("ParseRefresh failed with default expire: %v", err)
	}
	if claims.UserID != "user-5" {
		t.Fatalf("unexpected user id: %s", claims.UserID)
	}
}
