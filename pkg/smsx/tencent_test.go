package smsx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTencentSenderSendsRegisterCode(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("X-TC-Action") != "SendSms" {
			t.Fatalf("unexpected action: %s", r.Header.Get("X-TC-Action"))
		}
		if r.Header.Get("X-TC-Version") != "2021-01-11" {
			t.Fatalf("unexpected version: %s", r.Header.Get("X-TC-Version"))
		}
		if r.Header.Get("Authorization") == "" {
			t.Fatal("authorization header is required")
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_, _ = w.Write([]byte(`{"Response":{"SendStatusSet":[{"Code":"Ok","Message":"send success","PhoneNumber":"+8613800138000"}],"RequestId":"request-id"}}`))
	}))
	defer server.Close()

	sender := NewTencentSender(Config{
		SecretID:   "secret-id",
		SecretKey:  "secret-key",
		SDKAppID:   "1400000000",
		SignName:   "旅行助手",
		TemplateID: "123456",
		Region:     "ap-guangzhou",
		Endpoint:   server.URL,
		HTTPClient: &http.Client{Timeout: time.Second},
	})

	if err := sender.SendRegisterCode(context.Background(), "13800138000", "246810", 5*time.Minute); err != nil {
		t.Fatalf("send code failed: %v", err)
	}
	if got["SmsSdkAppId"] != "1400000000" || got["SignName"] != "旅行助手" || got["TemplateId"] != "123456" {
		t.Fatalf("unexpected sms payload: %#v", got)
	}
	params := got["TemplateParamSet"].([]any)
	if len(params) != 2 || params[0] != "246810" || params[1] != "5" {
		t.Fatalf("unexpected template params: %#v", params)
	}
	phones := got["PhoneNumberSet"].([]any)
	if len(phones) != 1 || phones[0] != "+8613800138000" {
		t.Fatalf("unexpected phone set: %#v", phones)
	}
}
