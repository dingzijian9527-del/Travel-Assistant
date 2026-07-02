package handler

import (
	"context"
	"testing"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/verifycodex"
)

func TestValidateRegisterCodeRejectsEmptyCode(t *testing.T) {
	store := verifycodex.NewMemoryStore()

	err := validateRegisterCode(context.Background(), store, "13800138000", "")
	if err == nil {
		t.Fatal("空验证码必须被拒绝")
	}
}

func TestValidateRegisterCodeRejectsWrongCode(t *testing.T) {
	store := verifycodex.NewMemoryStore()
	if err := store.Save(context.Background(), "13800138000", "246810", time.Minute); err != nil {
		t.Fatalf("保存验证码失败: %v", err)
	}

	err := validateRegisterCode(context.Background(), store, "13800138000", "135790")
	if err == nil {
		t.Fatal("错误验证码必须被拒绝")
	}
}

func TestValidateRegisterCodeAcceptsCorrectCode(t *testing.T) {
	store := verifycodex.NewMemoryStore()
	if err := store.Save(context.Background(), "13800138000", "246810", time.Minute); err != nil {
		t.Fatalf("保存验证码失败: %v", err)
	}

	if err := validateRegisterCode(context.Background(), store, "13800138000", "246810"); err != nil {
		t.Fatalf("正确验证码不应被拒绝: %v", err)
	}
	if err := validateRegisterCode(context.Background(), store, "13800138000", "246810"); err != nil {
		t.Fatalf("注册前校验不应提前删除验证码: %v", err)
	}
}
