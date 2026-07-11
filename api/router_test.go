package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestRegisterIncludesAIAgentStreamRoute(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App:  config.AppConfig{Name: "test"},
			AI:   config.AIConfig{MaxPromptChars: 2000},
			Auth: config.AuthConfig{JWTSecret: "unit-test-secret"},
			Redis: config.RedisConfig{Addr: "127.0.0.1:6379"},
		},
		Logger: zap.NewNop(),
	})

	for _, route := range h.Routes() {
		if route.Method == "POST" && route.Path == "/api/v1/ai-stream" {
			return
		}
	}
	t.Fatalf("ai stream route is not registered: %#v", h.Routes())
}

func TestRegisterIncludesUserRoutes(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App:  config.AppConfig{Name: "test"},
			AI:   config.AIConfig{MaxPromptChars: 2000},
			Auth: config.AuthConfig{JWTSecret: "unit-test-secret"},
			Redis: config.RedisConfig{Addr: "127.0.0.1:6379"},
		},
		Logger: zap.NewNop(),
	})

	requiredRoutes := map[string]bool{
		"POST /api/v1/user/register":     false,
		"POST /api/v1/user/login":        false,
		"GET /api/v1/user/profile":       false,
		"POST /api/v1/user/profile":      false,
		"GET /api/v1/user/dashboard":     false,
		"GET /api/v1/user/preferences":   false,
		"POST /api/v1/user/preferences":  false,
		"GET /api/v1/user/settings":      false,
		"POST /api/v1/user/settings":     false,
		"POST /api/v1/sms/register-code": false,
	}
	for _, route := range h.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := requiredRoutes[key]; ok {
			requiredRoutes[key] = true
		}
	}
	for route, found := range requiredRoutes {
		if !found {
			t.Fatalf("user route %s is not registered: %#v", route, h.Routes())
		}
	}
}

func TestRegisterIncludesTripRoutes(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App:  config.AppConfig{Name: "test"},
			AI:   config.AIConfig{MaxPromptChars: 2000},
			Auth: config.AuthConfig{JWTSecret: "unit-test-secret"},
			Redis: config.RedisConfig{Addr: "127.0.0.1:6379"},
		},
		Logger: zap.NewNop(),
	})

	requiredRoutes := map[string]bool{
		"POST /api/v1/trips":           false,
		"GET /api/v1/trips":            false,
		"GET /api/v1/trips/latest":     false,
		"GET /api/v1/trips/:tripID":    false,
		"DELETE /api/v1/trips/:tripID": false,
	}
	for _, route := range h.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := requiredRoutes[key]; ok {
			requiredRoutes[key] = true
		}
	}
	for route, found := range requiredRoutes {
		if !found {
			t.Fatalf("trip route %s is not registered: %#v", route, h.Routes())
		}
	}
}

type requestHeader struct {
	Key   string
	Value string
}

func TestRegisterAppliesRateLimitToLoginRoute(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App:   config.AppConfig{Name: "test", Env: "dev"},
			AI:    config.AIConfig{MaxPromptChars: 2000},
			Auth:  config.AuthConfig{JWTSecret: "unit-test-secret"},
			Redis: config.RedisConfig{Addr: "127.0.0.1:6379"},
		},
		Logger: zap.NewNop(),
	})

	body := []byte(`{"phone":"13800138000","password":"safe-pass-123"}`)
	headers := []requestHeader{{Key: "Content-Type", Value: "application/json"}, {Key: "X-Forwarded-For", Value: "127.0.0.1"}}
	for i := 0; i < 5; i++ {
		_ = performGatewayRequest(h, consts.MethodPost, "/api/v1/user/login", body, headers...)
	}
	sixth := performGatewayRequest(h, consts.MethodPost, "/api/v1/user/login", body, headers...)
	if sixth.Code != consts.StatusTooManyRequests {
		t.Fatalf("sixth login status = %d, want %d", sixth.Code, consts.StatusTooManyRequests)
	}
}

func TestRegisterAppliesRateLimitToSMSRoute(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App:   config.AppConfig{Name: "test", Env: "dev"},
			AI:    config.AIConfig{MaxPromptChars: 2000},
			Auth:  config.AuthConfig{JWTSecret: "unit-test-secret"},
			Redis: config.RedisConfig{Addr: "127.0.0.1:6379"},
			SMS:   config.SMSConfig{RegisterCodeExpire: 5},
		},
		Logger: zap.NewNop(),
	})

	body := []byte(`{"phone":"13800138000"}`)
	headers := []requestHeader{{Key: "Content-Type", Value: "application/json"}, {Key: "X-Forwarded-For", Value: "127.0.0.1"}, {Key: "X-Device-ID", Value: "device-1"}}
	_ = performGatewayRequest(h, consts.MethodPost, "/api/v1/sms/register-code", body, headers...)
	second := performGatewayRequest(h, consts.MethodPost, "/api/v1/sms/register-code", body, headers...)
	if second.Code != consts.StatusTooManyRequests {
		t.Fatalf("second sms status = %d, want %d", second.Code, consts.StatusTooManyRequests)
	}
}

func performGatewayRequest(h *server.Hertz, method string, path string, body []byte, headers ...requestHeader) *httptest.ResponseRecorder {
	ctx := h.Engine.NewContext()
	request := protocol.NewRequest(method, path, nil)
	request.CopyTo(&ctx.Request)
	if body != nil {
		ctx.Request.SetBody(body)
	}
	for _, item := range headers {
		ctx.Request.Header.Set(item.Key, item.Value)
	}

	h.ServeHTTP(context.Background(), ctx)

	writer := httptest.NewRecorder()
	ctx.Response.Header.VisitAll(func(key, value []byte) {
		writer.Header().Add(string(key), string(value))
	})
	writer.WriteHeader(ctx.Response.StatusCode())
	_, _ = writer.Write(ctx.Response.Body())
	ctx.Reset()
	return writer
}
