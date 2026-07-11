package handler

import (
	"strings"
	"testing"
)

func TestValidatePasswordRejectsTooShort(t *testing.T) {
	err := validatePassword("Ab1")
	if err == nil {
		t.Fatal("3位密码应被拒绝")
	}
	if !strings.Contains(err.Error(), "8位") {
		t.Fatalf("期望提示长度不足, 实际: %v", err)
	}
}

func TestValidatePasswordRejectsTooLong(t *testing.T) {
	long := strings.Repeat("Abc12345", 10) // 80 chars
	err := validatePassword(long)
	if err == nil {
		t.Fatal("超64位密码应被拒绝")
	}
}

func TestValidatePasswordRejectsMissingDigit(t *testing.T) {
	err := validatePassword("Abcdefgh")
	if err == nil {
		t.Fatal("缺少数字应被拒绝")
	}
	if !strings.Contains(err.Error(), "数字") {
		t.Fatalf("期望提示缺少数字, 实际: %v", err)
	}
}

func TestValidatePasswordRejectsMissingLower(t *testing.T) {
	err := validatePassword("ABC12345")
	if err == nil {
		t.Fatal("缺少小写字母应被拒绝")
	}
}

func TestValidatePasswordRejectsMissingUpper(t *testing.T) {
	err := validatePassword("abc12345")
	if err == nil {
		t.Fatal("缺少大写字母应被拒绝")
	}
}

func TestValidatePasswordRejectsEmpty(t *testing.T) {
	err := validatePassword("")
	if err == nil {
		t.Fatal("空密码应被拒绝")
	}
}

func TestValidatePasswordAcceptsValid(t *testing.T) {
	for _, pwd := range []string{
		"Abc12345",
		"MyP@ssw0rd2024",
		"Zzzzzzz1aB",
		"HelloWorld8",
	} {
		if err := validatePassword(pwd); err != nil {
			t.Fatalf("密码 %q 应通过但被拒绝: %v", pwd, err)
		}
	}
}
