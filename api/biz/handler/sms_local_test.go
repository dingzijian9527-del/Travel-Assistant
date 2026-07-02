package handler

import (
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestShouldUseLocalRegisterCodeWhenDevAndTencentMissing(t *testing.T) {
	runtime := &bootstrap.Runtime{Config: &config.Config{
		App: config.AppConfig{Env: "dev"},
		SMS: config.SMSConfig{DevReturnCode: true},
	}}

	if !shouldUseLocalRegisterCode(runtime) {
		t.Fatal("开发环境且腾讯云短信未配置时，应允许返回本地随机验证码")
	}
}

func TestShouldNotUseLocalRegisterCodeWhenTencentConfigured(t *testing.T) {
	runtime := &bootstrap.Runtime{Config: &config.Config{
		App: config.AppConfig{Env: "dev"},
		SMS: config.SMSConfig{
			DevReturnCode: true,
			SecretID:      "secret-id",
			SecretKey:     "secret-key",
			SDKAppID:      "1400000000",
			SignName:      "旅行助手",
			TemplateID:    "123456",
		},
	}}

	if shouldUseLocalRegisterCode(runtime) {
		t.Fatal("腾讯云短信配置完整时，应走真实短信发送")
	}
}

func TestShouldNotUseLocalRegisterCodeInProduction(t *testing.T) {
	runtime := &bootstrap.Runtime{Config: &config.Config{
		App: config.AppConfig{Env: "prod"},
		SMS: config.SMSConfig{DevReturnCode: true},
	}}

	if shouldUseLocalRegisterCode(runtime) {
		t.Fatal("生产环境禁止返回明文验证码")
	}
}
