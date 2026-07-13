package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func TestMemoryConfirmTokenStoreSaveAndVerify(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	ctx := context.Background()

	token := "test-confirm-token-12345"
	if err := store.Save(ctx, token, time.Minute); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	ok, err := store.Verify(ctx, token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if !ok {
		t.Fatal("Verify should return true for valid token")
	}
}

func TestMemoryConfirmTokenStoreRejectsWrongToken(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	ctx := context.Background()

	if err := store.Save(ctx, "right-token", time.Minute); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	ok, err := store.Verify(ctx, "wrong-token")
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if ok {
		t.Fatal("Verify should return false for wrong token")
	}
}

func TestMemoryConfirmTokenStoreRejectsExpiredToken(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	ctx := context.Background()

	if err := store.Save(ctx, "expired-token", 1*time.Millisecond); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	ok, err := store.Verify(ctx, "expired-token")
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if ok {
		t.Fatal("Verify should return false for expired token")
	}
}

func TestMemoryConfirmTokenStoreDelete(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	ctx := context.Background()

	if err := store.Save(ctx, "delete-me", time.Minute); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := store.Delete(ctx, "delete-me"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	ok, err := store.Verify(ctx, "delete-me")
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if ok {
		t.Fatal("Verify should return false after delete")
	}
}

func TestGenerateConfirmToken(t *testing.T) {
	token1, err := GenerateConfirmToken()
	if err != nil {
		t.Fatalf("GenerateConfirmToken failed: %v", err)
	}
	if len(token1) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(token1))
	}

	token2, err := GenerateConfirmToken()
	if err != nil {
		t.Fatalf("second GenerateConfirmToken failed: %v", err)
	}
	if token1 == token2 {
		t.Fatal("two generated tokens should be different")
	}
}

func TestRequirePasswordConfirmRejectsMissingHeader(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	handler := RequirePasswordConfirm(store)

	c := &app.RequestContext{}
	c.Request.SetHost("localhost:8080")
	handler(context.Background(), c)

	if c.Response.StatusCode() != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", c.Response.StatusCode())
	}
}

func TestRequirePasswordConfirmAcceptsValidToken(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	ctx := context.Background()

	token := "valid-token-abc"
	if err := store.Save(ctx, token, time.Minute); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	handler := RequirePasswordConfirm(store)

	c := &app.RequestContext{}
	c.Request.SetHost("localhost:8080")
	c.Request.Header.Set("X-Confirm-Token", token)
	handler(context.Background(), c)

	if c.Response.StatusCode() != http.StatusOK {
		t.Fatalf("expected 200, got %d", c.Response.StatusCode())
	}

	// Token should be consumed (one-time use)
	ok, _ := store.Verify(ctx, token)
	if ok {
		t.Fatal("token should be deleted after use")
	}
}

func TestRequirePasswordConfirmRejectsInvalidToken(t *testing.T) {
	store := NewMemoryConfirmTokenStore()
	handler := RequirePasswordConfirm(store)

	c := &app.RequestContext{}
	c.Request.SetHost("localhost:8080")
	c.Request.Header.Set("X-Confirm-Token", "invalid-token")
	handler(context.Background(), c)

	if c.Response.StatusCode() != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", c.Response.StatusCode())
	}
}
