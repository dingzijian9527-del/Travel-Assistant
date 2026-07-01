package rpcaiagent

import (
	"bytes"
	"context"
	"errors"
	"strings"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	travelai "github.com/dingzijian9527-del/Travel-Assistant/pkg/ai"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

type aiAgentService struct {
	repo  *aiAgentRepo
	model modelClient
}

func newAIAgentService(repo *aiAgentRepo) *aiAgentService {
	return newAIAgentServiceWithModel(repo, newTravelModelClient(config.AIConfig{}))
}

func newAIAgentServiceWithModel(repo *aiAgentRepo, model modelClient) *aiAgentService {
	if model == nil {
		model = newTravelModelClient(config.AIConfig{})
	}
	return &aiAgentService{repo: repo, model: model}
}

func (l *aiAgentService) Chat(ctx context.Context, req *aiagent.ChatRequest) (*chatMessageModel, []string, *serviceError) {
	message := strings.TrimSpace(req.Message)
	if req.UserId <= 0 {
		return nil, nil, errParam("valid userId is required")
	}
	if message == "" {
		return nil, nil, errParam("message is required")
	}
	if len([]rune(message)) > 1000 {
		return nil, nil, errParam("message is too long")
	}
	userMessage := &chatMessageModel{Role: aiagent.ChatRole_USER, Content: message}
	history := l.repo.History(ctx, req.UserId)
	if !isTravelRelated(message) && len(history) == 0 {
		reply := "您好，我是旅行助手的旅游专用智能体，主要帮助你规划行程、制作攻略、推荐住宿交通和解答目的地相关问题。这个问题不属于旅游出行场景，我暂时不能展开回答。"
		assistantMessage := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: reply}
		l.repo.AppendMessages(ctx, req.UserId, userMessage, assistantMessage)
		return assistantMessage, buildSuggestions(message), nil
	}
	reply, err := l.model.Chat(ctx, buildModelMessages(history, message))
	if err != nil {
		return nil, nil, errParam("智能体模型暂不可用，请稍后重试")
	}
	assistantMessage := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: reply}
	l.repo.AppendMessages(ctx, req.UserId, userMessage, assistantMessage)
	return assistantMessage, buildSuggestions(message), nil
}

func (l *aiAgentService) GetPromptSuggestions(_ context.Context, req *aiagent.PromptSuggestionsRequest) []string {
	scene := ""
	if req.Scene != nil {
		scene = strings.TrimSpace(*req.Scene)
	}
	switch scene {
	case "hotel":
		return []string{"搜索上海精品酒店", "推荐三亚亲子酒店", "帮我比较民宿和酒店"}
	case "plan":
		return []string{"定制 5 天大理行程", "规划北京亲子 3 日游", "生成三亚海岛度假路线"}
	case "guide":
		return []string{"制作西安旅游攻略", "推荐成都本地美食", "整理杭州避坑指南"}
	default:
		return []string{"定制 5 天大理行程", "搜索上海精品酒店", "查看今日天气"}
	}
}

type modelMessage struct {
	Role    string
	Content string
}

type modelClient interface {
	Chat(ctx context.Context, messages []modelMessage) (string, error)
}

type travelModelClient struct {
	client *travelai.Client
}

func newTravelModelClient(cfg config.AIConfig) *travelModelClient {
	return &travelModelClient{client: travelai.NewClient(cfg)}
}

func (c *travelModelClient) Chat(ctx context.Context, messages []modelMessage) (string, error) {
	if c.client == nil || !c.client.Available() {
		return "", errors.New("model unavailable")
	}
	converted := make([]travelai.Message, 0, len(messages))
	for _, item := range messages {
		converted = append(converted, travelai.Message{Role: item.Role, Content: item.Content})
	}
	var output bytes.Buffer
	if err := c.client.StreamChatWithMessages(ctx, converted, &output); err != nil {
		return "", err
	}
	return strings.TrimSpace(output.String()), nil
}

func buildModelMessages(history []chatMessageModel, message string) []modelMessage {
	messages := make([]modelMessage, 0, len(history)+1)
	for _, item := range history {
		role := "user"
		if item.Role == aiagent.ChatRole_ASSISTANT {
			role = "assistant"
		}
		if strings.TrimSpace(item.Content) != "" {
			messages = append(messages, modelMessage{Role: role, Content: item.Content})
		}
	}
	messages = append(messages, modelMessage{Role: "user", Content: message})
	return messages
}
