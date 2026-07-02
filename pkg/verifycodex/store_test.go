package verifycodex

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStoreSaveAndVerifyCode(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	if err := store.Save(ctx, "13800138000", "246810", time.Minute); err != nil {
		t.Fatalf("save code failed: %v", err)
	}
	ok, err := store.Verify(ctx, "13800138000", "246810")
	if err != nil {
		t.Fatalf("verify code failed: %v", err)
	}
	if !ok {
		t.Fatal("expected code to match")
	}
	ok, err = store.Verify(ctx, "13800138000", "246810")
	if err != nil {
		t.Fatalf("verify deleted code failed: %v", err)
	}
	if ok {
		t.Fatal("code should be deleted after successful verification")
	}
}

func TestMemoryStoreCheckDoesNotDeleteCode(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	if err := store.Save(ctx, "13800138000", "246810", time.Minute); err != nil {
		t.Fatalf("save code failed: %v", err)
	}
	ok, err := store.Check(ctx, "13800138000", "246810")
	if err != nil {
		t.Fatalf("check code failed: %v", err)
	}
	if !ok {
		t.Fatal("expected code to match")
	}
	ok, err = store.Check(ctx, "13800138000", "246810")
	if err != nil {
		t.Fatalf("check code again failed: %v", err)
	}
	if !ok {
		t.Fatal("check should not delete code")
	}
}

func TestMemoryStoreRejectsExpiredCode(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	if err := store.Save(ctx, "13800138000", "246810", -time.Second); err != nil {
		t.Fatalf("save code failed: %v", err)
	}
	ok, err := store.Verify(ctx, "13800138000", "246810")
	if err != nil {
		t.Fatalf("verify code failed: %v", err)
	}
	if ok {
		t.Fatal("expired code should not match")
	}
}
