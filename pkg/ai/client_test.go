package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestClientStreamsChatCompletionWithHistory(t *testing.T) {
	var authHeader string
	var requestBody chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"上海\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"天气建议\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	client := NewClient(config.AIConfig{
		APIKey:       "test-key",
		BaseURL:      server.URL,
		EndpointID:   "ep-test",
		ModelName:    "DeepSeek-V4-pro",
		Timeout:      time.Second,
		Stream:       true,
		SystemPrompt: DefaultSystemPrompt(),
	})
	var output strings.Builder
	err := client.StreamChatWithMessages(context.Background(), []Message{
		{Role: "user", Content: "我要去上海玩7天"},
		{Role: "assistant", Content: "上海7天路线建议"},
		{Role: "user", Content: "天气如何"},
	}, &output)
	if err != nil {
		t.Fatalf("stream chat returned error: %v", err)
	}

	if authHeader != "Bearer test-key" {
		t.Fatalf("authorization header mismatch: %s", authHeader)
	}
	if requestBody.Model != "ep-test" || !requestBody.Stream {
		t.Fatalf("request model or stream mismatch: %#v", requestBody)
	}
	if len(requestBody.Messages) != 4 {
		t.Fatalf("expected system plus history messages, got %#v", requestBody.Messages)
	}
	if requestBody.Messages[1].Role != "user" || requestBody.Messages[2].Role != "assistant" || requestBody.Messages[3].Content != "天气如何" {
		t.Fatalf("history messages not preserved: %#v", requestBody.Messages)
	}
	if output.String() != "上海天气建议" {
		t.Fatalf("stream output mismatch: %s", output.String())
	}
}

func TestClientUnavailableWithoutEndpointOrModel(t *testing.T) {
	client := NewClient(config.AIConfig{
		APIKey:  "test-key",
		BaseURL: "https://example.test",
	})
	if client.Available() {
		t.Fatal("client should be unavailable without endpoint or model")
	}
}

func TestClientFallsBackToLegacyModelWhenEndpointMissing(t *testing.T) {
	var requestBody chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	client := NewClient(config.AIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Model:   "legacy-model",
		Timeout: time.Second,
	})
	var output strings.Builder
	if err := client.StreamChat(context.Background(), "测试", &output); err != nil {
		t.Fatalf("stream chat returned error: %v", err)
	}
	if requestBody.Model != "legacy-model" {
		t.Fatalf("legacy model should be used when endpoint is missing: %#v", requestBody)
	}
}
