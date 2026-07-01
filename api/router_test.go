package api

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestRegisterIncludesAIAgentStreamRoute(t *testing.T) {
	h := server.New()
	Register(h, &bootstrap.Runtime{
		Config: &config.Config{
			App: config.AppConfig{Name: "test"},
			AI:  config.AIConfig{MaxPromptChars: 2000},
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
			App: config.AppConfig{Name: "test"},
			AI:  config.AIConfig{MaxPromptChars: 2000},
		},
		Logger: zap.NewNop(),
	})

	requiredRoutes := map[string]bool{
		"POST /api/v1/user/register": false,
		"POST /api/v1/user/login":    false,
		"GET /api/v1/user/profile":   false,
		"POST /api/v1/user/profile":  false,
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
