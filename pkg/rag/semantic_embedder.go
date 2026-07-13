package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

type ContextEmbedder interface {
	EmbedContext(ctx context.Context, text string) ([]float32, error)
}

type SemanticEmbedder struct {
	endpoint string
	apiKey   string
	model    string
	dim      int
	client   *http.Client
}

type semanticEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type semanticEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func NewConfiguredEmbedder(cfg config.RAGConfig) (Embedder, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "", "local":
		return NewHashEmbedder(cfg.EmbeddingDim), nil
	case "semantic", "openai", "ark":
		return NewSemanticEmbedder(cfg)
	default:
		return nil, fmt.Errorf("不支持的嵌入提供方：%s", cfg.Provider)
	}
}

func NewSemanticEmbedder(cfg config.RAGConfig) (*SemanticEmbedder, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(cfg.EmbeddingBaseURL), "/")
	if endpoint == "" {
		return nil, errors.New("语义嵌入服务地址不能为空")
	}
	if !strings.HasSuffix(endpoint, "/embeddings") {
		endpoint += "/embeddings"
	}
	if strings.TrimSpace(cfg.EmbeddingAPIKey) == "" {
		return nil, errors.New("语义嵌入服务密钥不能为空")
	}
	if strings.TrimSpace(cfg.EmbeddingModel) == "" {
		return nil, errors.New("语义嵌入模型不能为空")
	}
	return &SemanticEmbedder{
		endpoint: endpoint,
		apiKey:   strings.TrimSpace(cfg.EmbeddingAPIKey),
		model:    strings.TrimSpace(cfg.EmbeddingModel),
		dim:      cfg.EmbeddingDim,
		client:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (e *SemanticEmbedder) Embed(text string) []float32 {
	vector, _ := e.EmbedContext(context.Background(), text)
	return vector
}

func (e *SemanticEmbedder) EmbedContext(ctx context.Context, text string) ([]float32, error) {
	payload, err := json.Marshal(semanticEmbeddingRequest{Model: e.model, Input: text})
	if err != nil {
		return nil, fmt.Errorf("编码嵌入请求失败：%w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("创建嵌入请求失败：%w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+e.apiKey)
	response, err := e.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("调用嵌入服务失败：%w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("嵌入服务返回状态异常：%d", response.StatusCode)
	}
	var result semanticEmbeddingResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析嵌入响应失败：%w", err)
	}
	if len(result.Data) == 0 || len(result.Data[0].Embedding) == 0 {
		return nil, errors.New("嵌入响应没有向量")
	}
	vector := result.Data[0].Embedding
	if e.dim > 0 && len(vector) != e.dim {
		return nil, fmt.Errorf("嵌入维度不匹配：期望 %d，实际 %d", e.dim, len(vector))
	}
	return vector, nil
}
