package middleware

import (
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
)

func TestHashLoginFingerprint_DifferentIPs(t *testing.T) {
	c1 := &app.RequestContext{}
	c1.Request.SetHost("10.0.0.1:8080")
	c1.Request.Header.Set("X-Forwarded-For", "1.2.3.4")
	c2 := &app.RequestContext{}
	c2.Request.SetHost("10.0.0.1:8080")
	c2.Request.Header.Set("X-Forwarded-For", "5.6.7.8")

	fp1 := hashLoginFingerprint("13800138000", c1)
	fp2 := hashLoginFingerprint("13800138000", c2)

	if fp1 == fp2 {
		t.Fatalf("同一个手机号不同 IP 应生成不同指纹")
	}
}

func TestHashLoginFingerprint_SameInput(t *testing.T) {
	c1 := &app.RequestContext{}
	c1.Request.SetHost("10.0.0.1:8080")
	c2 := &app.RequestContext{}
	c2.Request.SetHost("10.0.0.1:8080")

	fp1 := hashLoginFingerprint("13800138000", c1)
	fp2 := hashLoginFingerprint("13800138000", c2)

	if fp1 != fp2 {
		t.Fatalf("相同手机号和 IP 应生成相同指纹")
	}
}

func TestHashLoginFingerprint_UsesXForwardedFor(t *testing.T) {
	c := &app.RequestContext{}
	c.Request.SetHost("10.0.0.1:8080")
	c.Request.Header.Set("X-Forwarded-For", "203.0.113.5")

	fp := hashLoginFingerprint("13800138000", c)
	if fp == "" {
		t.Fatal("指纹不应为空")
	}
}

func TestLoginPhoneFromRequest_ExtractsPhone(t *testing.T) {
	body, _ := json.Marshal(map[string]any{
		"phone":    "13800138000",
		"password": "test123",
	})

	c := &app.RequestContext{}
	c.Request.SetRequestURI("/api/v1/user/login")
	c.Request.SetBody(body)

	phone := loginPhoneFromRequest(c)
	if phone != "13800138000" {
		t.Fatalf("expected 13800138000, got %q", phone)
	}
}

func TestLoginPhoneFromRequest_SkipsNonLoginPaths(t *testing.T) {
	body, _ := json.Marshal(map[string]any{"phone": "13800138000"})

	c := &app.RequestContext{}
	c.Request.SetRequestURI("/api/v1/sms/register-code")
	c.Request.SetBody(body)

	phone := loginPhoneFromRequest(c)
	if phone != "" {
		t.Fatalf("非登录接口不应提取手机号, got %q", phone)
	}
}

func TestLoginPhoneFromRequest_EmptyBody(t *testing.T) {
	c := &app.RequestContext{}
	c.Request.SetRequestURI("/api/v1/user/login")

	phone := loginPhoneFromRequest(c)
	if phone != "" {
		t.Fatalf("空 body 应返回空, got %q", phone)
	}
}
