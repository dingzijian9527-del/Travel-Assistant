package middleware

import (
	"testing"
	"time"
)

func TestLoginFailureLimiterBlocksAnyRepeatedIdentityDimension(t *testing.T) {
	limiter := NewLoginFailureLimiter(2, time.Minute)
	first := loginIdentity{Phone: "13800138000", IP: "10.0.0.1", DeviceID: "device-a"}
	secondPhone := loginIdentity{Phone: "13800138001", IP: "10.0.0.1", DeviceID: "device-b"}

	limiter.RecordFailure(first)
	limiter.RecordFailure(first)
	if limiter.Allow(first) {
		t.Fatal("同一身份连续失败后必须拦截")
	}
	if limiter.Allow(secondPhone) {
		t.Fatal("更换手机号但复用来源地址后仍必须拦截")
	}
}

func TestLoginFailureLimiterResetsAfterSuccessfulLogin(t *testing.T) {
	limiter := NewLoginFailureLimiter(2, time.Minute)
	identity := loginIdentity{Phone: "13800138000", IP: "10.0.0.1", DeviceID: "device-a"}

	limiter.RecordFailure(identity)
	limiter.RecordFailure(identity)
	limiter.Reset(identity)
	if !limiter.Allow(identity) {
		t.Fatal("登录成功后应清除失败计数")
	}
}
