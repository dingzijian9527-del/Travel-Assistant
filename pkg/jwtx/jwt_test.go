package jwtx

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	token, err := Generate(Config{Secret: "test-secret", Expire: time.Hour}, Claims{UserID: "user-1", Phone: "13800138000"})
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	claims, err := Parse(Config{Secret: "test-secret", Expire: time.Hour}, token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.UserID != "user-1" || claims.Phone != "13800138000" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseRejectsTamperedToken(t *testing.T) {
	token, err := Generate(Config{Secret: "test-secret", Expire: time.Hour}, Claims{UserID: "user-1"})
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	if _, err := Parse(Config{Secret: "test-secret", Expire: time.Hour}, token+"x"); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}
