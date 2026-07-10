package rpcaiagent

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	travelai "github.com/dingzijian9527-del/Travel-Assistant/pkg/ai"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/traveldata"
)

// aiAgentService 实现 AI 智能体的核心业务逻辑。
// 包含 RAG 检索增强、对话上下文管理、流式与非流式对话。
type aiAgentService struct {
	repo          aiAgentRepository
	model         modelClient
	retriever     rag.Retriever
	ragCfg        config.RAGConfig
	travelPlanner travelDataPlanner
}

type travelDataPlanner interface {
	BuildContext(ctx context.Context, req traveldata.Request) (traveldata.Result, error)
}

// newAIAgentService 创建默认的 AI 智能体服务（用于测试）。
func newAIAgentService(repo aiAgentRepository) *aiAgentService {
	return newAIAgentServiceWithModel(repo, newTravelModelClient(config.AIConfig{}), nil, config.RAGConfig{})
}

// newAIAgentServiceWithModel 创建指定模型和检索器的 AI 智能体服务。
func newAIAgentServiceWithModel(repo aiAgentRepository, model modelClient, retriever rag.Retriever, ragCfg config.RAGConfig) *aiAgentService {
	return newAIAgentServiceWithTravelData(repo, model, retriever, ragCfg, nil)
}

func newAIAgentServiceWithTravelData(repo aiAgentRepository, model modelClient, retriever rag.Retriever, ragCfg config.RAGConfig, travelPlanner travelDataPlanner) *aiAgentService {
	if model == nil {
		model = newTravelModelClient(config.AIConfig{})
	}
	return &aiAgentService{repo: repo, model: model, retriever: retriever, ragCfg: ragCfg, travelPlanner: travelPlanner}
}

// Chat 实现单轮非流式对话，包含 RAG 检索增强。
// 流程：校验 -> 检索知识 -> 构建提示词 -> 调用大模型 -> 返回回复
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

	// RAG 检索旅行知识
	referenceContext := l.buildReferenceContext(ctx, message)

	if !isTravelRelated(message) && len(history) == 0 && referenceContext == "" {
		reply := "你好，我是旅行助手的旅游专用智能体，主要帮助你规划行程、制作攻略、推荐住宿交通和解答目的地相关问题。这个问题不属于旅游出行场景，我暂时不能展开回答。"
		assistantMessage := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: reply}
		l.repo.AppendMessages(ctx, req.UserId, userMessage, assistantMessage)
		return assistantMessage, buildSuggestions(message), nil
	}

	// 构建带 RAG 上下文的模型消息
	modelMessages := buildModelMessages(history, message, referenceContext)
	reply, err := l.model.Chat(ctx, modelMessages)
	if err != nil {
		return nil, nil, errParam("智能体模型暂不可用，请稍后重试")
	}

	assistantMessage := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: reply}
	l.repo.AppendMessages(ctx, req.UserId, userMessage, assistantMessage)
	return assistantMessage, buildSuggestions(message), nil
}

// ChatStream 实现流式对话，包含 RAG 检索增强。
func (l *aiAgentService) ChatStream(ctx context.Context, userID int64, message string, replyWriter io.Writer) *serviceError {
	message = strings.TrimSpace(message)
	if message == "" {
		return errParam("message is required")
	}
	if len([]rune(message)) > 1000 {
		return errParam("message is too long")
	}

	userMsg := &chatMessageModel{Role: aiagent.ChatRole_USER, Content: message}
	history := l.repo.History(ctx, userID)

	// RAG 检索旅行知识
	referenceContext := l.buildReferenceContext(ctx, message)

	// 非旅游问题且无历史记录/知识，引导提示
	if !isTravelRelated(message) && len(history) == 0 && referenceContext == "" {
		guideReply := "你好，我是旅行助手的旅游专用智能体，主要帮助你规划行程、制作攻略、推荐住宿交通和解答目的地相关问题。这个问题不属于旅游出行场景，我暂时不能展开回答。"
		if _, err := replyWriter.Write([]byte(guideReply)); err != nil {
			return errParam("写入回复失败")
		}
		assistantMsg := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: guideReply}
		l.repo.AppendMessages(ctx, userID, userMsg, assistantMsg)
		return nil
	}

	// 构建带 RAG 上下文的模型消息
	modelMessages := buildModelMessages(history, message, referenceContext)

	var replyBuffer bytes.Buffer
	streamWriter := io.MultiWriter(replyWriter, &replyBuffer)

	modelUnavailableErr := l.model.StreamChat(ctx, modelMessages, streamWriter)
	if modelUnavailableErr != nil {
		if replyBuffer.Len() == 0 {
			if _, writeErr := replyWriter.Write([]byte("智能体模型暂不可用，请稍后重试")); writeErr != nil {
				return errParam("写入回复失败")
			}
		}
		return nil
	}

	// 保存对话历史
	replyText := strings.TrimSpace(replyBuffer.String())
	if replyText != "" {
		assistantMsg := &chatMessageModel{Role: aiagent.ChatRole_ASSISTANT, Content: replyText}
		l.repo.AppendMessages(ctx, userID, userMsg, assistantMsg)
	}
	return nil
}

// retrieveKnowledge 执行 RAG 检索，返回格式化的参考资料文本。
func (l *aiAgentService) retrieveKnowledge(ctx context.Context, message string) string {
	if l.retriever == nil {
		return ""
	}
	results, err := l.retriever.Search(ctx, message, l.ragCfg)
	if err != nil {
		return ""
	}
	return rag.FormatContext(results)
}

func (l *aiAgentService) buildReferenceContext(ctx context.Context, message string) string {
	parts := make([]string, 0, 2)
	if knowledge := l.retrieveKnowledge(ctx, message); strings.TrimSpace(knowledge) != "" {
		parts = append(parts, knowledge)
	}
	if travelContext := l.retrieveTravelData(ctx, message); strings.TrimSpace(travelContext) != "" {
		parts = append(parts, travelContext)
	}
	return strings.Join(parts, "\n\n")
}

func (l *aiAgentService) retrieveTravelData(ctx context.Context, message string) string {
	if l.travelPlanner == nil || !isTravelRelated(message) {
		return ""
	}
	req := traveldata.Request{
		Message:     message,
		Destination: traveldata.InferDestination(message),
		Days:        traveldata.InferDays(message),
		People:      traveldata.InferPeople(message),
		Budget:      traveldata.InferBudget(message),
	}
	result, err := l.travelPlanner.BuildContext(ctx, req)
	if err != nil {
		return ""
	}
	return result.FormatForPrompt()
}

// GetPromptSuggestions 返回场景化的提示词推荐。
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

// modelClient 大模型客户端接口，支持流式和非流式调用。
type modelClient interface {
	Chat(ctx context.Context, messages []modelMessage) (string, error)
	StreamChat(ctx context.Context, messages []modelMessage, output io.Writer) error
}

// travelModelClient 旅行智能体的大模型客户端封装。
type travelModelClient struct {
	client *travelai.Client
}

// newTravelModelClient 创建大模型客户端。
func newTravelModelClient(cfg config.AIConfig) *travelModelClient {
	return &travelModelClient{client: travelai.NewClient(cfg)}
}

// Chat 非流式对话。
func (c *travelModelClient) Chat(ctx context.Context, messages []modelMessage) (string, error) {
	if c.client == nil || !c.client.Available() {
		return "", errParam("智能体模型暂不可用，请稍后重试")
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

// StreamChat 流式对话。
func (c *travelModelClient) StreamChat(ctx context.Context, messages []modelMessage, output io.Writer) error {
	if c.client == nil || !c.client.Available() {
		return errParam("智能体模型暂不可用，请稍后重试")
	}
	converted := make([]travelai.Message, 0, len(messages))
	for _, item := range messages {
		converted = append(converted, travelai.Message{Role: item.Role, Content: item.Content})
	}
	return c.client.StreamChatWithMessages(ctx, converted, output)
}

// modelMessage 表示发往大模型的一条消息。
type modelMessage struct {
	Role    string
	Content string
}

// buildModelMessages 构建发往大模型的消息列表，包含历史记录和 RAG 上下文。
func buildModelMessages(history []chatMessageModel, message string, referenceContext string) []modelMessage {
	messages := make([]modelMessage, 0, len(history)+2)
	messages = append(messages, modelMessage{Role: "user", Content: conversationStateRules()})

	for _, item := range history {
		role := "user"
		if item.Role == aiagent.ChatRole_ASSISTANT {
			role = "assistant"
		}
		if strings.TrimSpace(item.Content) != "" {
			messages = append(messages, modelMessage{Role: role, Content: item.Content})
		}
	}

	// 将 RAG 检索到的知识拼到当前用户消息前面
	currentMessage := message
	if strings.TrimSpace(referenceContext) != "" {
		currentMessage = referenceContext + "\n\n用户问题：" + message
	}
	messages = append(messages, modelMessage{Role: "user", Content: currentMessage})
	return messages
}

func conversationStateRules() string {
	return strings.TrimSpace(`
对话状态规则：
1. 必须结合最近对话历史理解用户当前问题，不要把每轮都当作新问题。
2. 如果历史里已经有目的地、时间、天数、预算、同行人或偏好，用户后续只补充“7月”“3天”“预算3000”“不要太赶”等片段时，必须沿用历史里已有条件。
3. 如果用户想要行程、天气、路线、住宿或推荐，但当前问题和历史都缺少具体地点或具体时间，先追问缺失信息，问题要简短明确。
4. 如果上一轮已经追问地点或时间，用户这一轮仍未提供缺失信息，不要继续反复追问，直接给出国内热门目的地、热门季节或通用周末方案，并说明用户补充地点或时间后可以继续细化。
5. 热门推荐要具体可执行，可以按亲子、美食、避暑、海岛、古城、周末短途等方向给出候选，并附上适合人群和注意点。
`)
}

// travelSessionStore 管理用户的对话历史（用于 API 网关层流式对话）。
type travelSessionStore struct {
	mu      sync.RWMutex
	max     int
	history map[int64][]conversationMessage
}

// conversationMessage 用于 API 网关层的简单对话消息。
type conversationMessage struct {
	Role    string
	Content string
}

// newTravelSessionStore 创建对话历史存储器。
func newTravelSessionStore(max int) *travelSessionStore {
	return &travelSessionStore{max: max, history: make(map[int64][]conversationMessage)}
}

// History 返回用户的对话历史副本。
func (s *travelSessionStore) History(userID int64) []conversationMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.history[userID]
	copied := make([]conversationMessage, len(items))
	copy(copied, items)
	return copied
}

// Append 追加一条用户消息和一条智能体回复到用户的历史记录中。
func (s *travelSessionStore) Append(userID int64, userText string, replyText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := append(s.history[userID],
		conversationMessage{Role: "user", Content: userText},
		conversationMessage{Role: "assistant", Content: replyText},
	)
	if len(items) > s.max {
		items = items[len(items)-s.max:]
	}
	s.history[userID] = items
}
