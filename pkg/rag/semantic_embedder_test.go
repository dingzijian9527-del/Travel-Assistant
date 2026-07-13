package rag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestSemanticEmbedderCallsConfiguredEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embeddings" {
			t.Fatalf("请求路径不正确：%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("鉴权请求头不正确")
		}
		var request struct {
			Model string `json:"model"`
			Input string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("请求体解析失败：%v", err)
		}
		if request.Model != "travel-embedding" || request.Input != "成都酒店推荐" {
			t.Fatalf("请求参数不正确：%+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"embedding": []float32{0.1, 0.2, 0.3}}},
		})
	}))
	defer server.Close()

	embedder, err := NewConfiguredEmbedder(config.RAGConfig{
		Provider:         "semantic",
		EmbeddingBaseURL: server.URL,
		EmbeddingAPIKey:  "test-key",
		EmbeddingModel:   "travel-embedding",
		EmbeddingDim:     3,
	})
	if err != nil {
		t.Fatalf("创建语义嵌入器失败：%v", err)
	}
	contextual, ok := embedder.(ContextEmbedder)
	if !ok {
		t.Fatal("语义嵌入器必须支持带上下文的调用")
	}
	vector, err := contextual.EmbedContext(context.Background(), "成都酒店推荐")
	if err != nil {
		t.Fatalf("调用语义嵌入接口失败：%v", err)
	}
	if len(vector) != 3 || vector[0] != 0.1 || vector[2] != 0.3 {
		t.Fatalf("向量结果不正确：%v", vector)
	}
}

func TestConfiguredEmbedderUsesHashOnlyForExplicitLocalProvider(t *testing.T) {
	embedder, err := NewConfiguredEmbedder(config.RAGConfig{Provider: "local", EmbeddingDim: 4})
	if err != nil {
		t.Fatalf("本地嵌入器创建失败：%v", err)
	}
	if _, ok := embedder.(*HashEmbedder); !ok {
		t.Fatalf("本地模式应使用明确的哈希嵌入器，实际类型为 %T", embedder)
	}
}
