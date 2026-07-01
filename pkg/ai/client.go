package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// Message 表示发送给人工智能接口的一条上下文消息。
type Message struct {
	Role    string
	Content string
}

// Client 封装兼容对话接口的旅行智能体调用能力。
type Client struct {
	cfg        config.AIConfig
	httpClient *http.Client
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Stream   bool          `json:"stream"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamResponse struct {
	Choices []streamChoice `json:"choices"`
}

type streamChoice struct {
	Delta streamDelta `json:"delta"`
}

type streamDelta struct {
	Content string `json:"content"`
}

// NewClient 创建旅行智能体接口客户端。
func NewClient(cfg config.AIConfig) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Available 判断当前配置是否足以调用远程智能体接口。
func (c *Client) Available() bool {
	return strings.TrimSpace(c.cfg.APIKey) != "" &&
		strings.TrimSpace(c.cfg.BaseURL) != "" &&
		c.modelIdentifier() != ""
}

// StreamChat 调用远程智能体接口，并将流式文本写入输出目标。
func (c *Client) StreamChat(ctx context.Context, userMessage string, output io.Writer) error {
	return c.StreamChatWithMessages(ctx, []Message{{Role: "user", Content: userMessage}}, output)
}

// StreamChatWithMessages 调用远程智能体接口，并携带最近对话上下文。
func (c *Client) StreamChatWithMessages(ctx context.Context, messages []Message, output io.Writer) error {
	if !c.Available() {
		return errors.New("智能体接口配置不完整")
	}
	requestBody := chatCompletionRequest{
		Model:    c.modelIdentifier(),
		Stream:   true,
		Messages: buildRequestMessages(c.cfg.SystemPrompt, messages),
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("编码智能体请求失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL(c.cfg.BaseURL), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("创建智能体请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.cfg.APIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("调用智能体接口失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("智能体接口状态异常: %d", resp.StatusCode)
	}
	return writeStreamContent(resp.Body, output)
}

func buildRequestMessages(prompt string, messages []Message) []chatMessage {
	result := []chatMessage{{Role: "system", Content: systemPrompt(prompt)}}
	for _, message := range messages {
		role := strings.TrimSpace(message.Role)
		content := strings.TrimSpace(message.Content)
		if role == "" || content == "" {
			continue
		}
		if role != "user" && role != "assistant" {
			role = "user"
		}
		result = append(result, chatMessage{Role: role, Content: content})
	}
	return result
}

func writeStreamContent(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return nil
		}
		var chunk streamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("解析智能体流式响应失败: %w", err)
		}
		for _, choice := range chunk.Choices {
			if choice.Delta.Content == "" {
				continue
			}
			if _, err := output.Write([]byte(choice.Delta.Content)); err != nil {
				return fmt.Errorf("写入智能体响应失败: %w", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取智能体流式响应失败: %w", err)
	}
	return nil
}

func chatCompletionsURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(trimmed, "/chat/completions") {
		return trimmed
	}
	return trimmed + "/chat/completions"
}

func (c *Client) modelIdentifier() string {
	if endpointID := strings.TrimSpace(c.cfg.EndpointID); endpointID != "" {
		return endpointID
	}
	return strings.TrimSpace(c.cfg.Model)
}

func systemPrompt(prompt string) string {
	if strings.TrimSpace(prompt) != "" {
		return strings.TrimSpace(prompt)
	}
	return DefaultSystemPrompt()
}
