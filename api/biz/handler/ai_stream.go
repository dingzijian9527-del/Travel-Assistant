package handler

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/validator"
	travelai "github.com/dingzijian9527-del/Travel-Assistant/pkg/ai"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
)

const modelUnavailableReply = "智能体模型暂不可用，请稍后重试或检查火山方舟接入点配置。"

var travelSessions = newTravelSessionStore(12)

// ChatStreamRequest 描述旅行智能体流式对话请求。
type ChatStreamRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type conversationMessage struct {
	Role    string
	Content string
}

type travelSessionStore struct {
	mu      sync.RWMutex
	max     int
	history map[int64][]conversationMessage
}

func newTravelSessionStore(max int) *travelSessionStore {
	return &travelSessionStore{max: max, history: make(map[int64][]conversationMessage)}
}

func (s *travelSessionStore) History(userID int64) []conversationMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.history[userID]
	copied := make([]conversationMessage, len(items))
	copy(copied, items)
	return copied
}

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

// ChatStream 处理旅行智能体流式对话。
func ChatStream(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req ChatStreamRequest
		if err := c.BindAndValidate(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": "请求参数错误"})
			return
		}
		message, ok, msg := validator.RequiredString(req.Message, "message")
		if !ok {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": msg})
			return
		}
		if utf8.RuneCountInString(message) > runtime.Config.AI.MaxPromptChars {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": "问题内容过长"})
			return
		}

		claims, ok := claimsFromRequest(runtime, c)
		if !ok {
			return
		}
		userID, err := jwtx.UserIDInt64(claims.UserID)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]any{"code": 401, "msg": "登录状态无效，请重新登录"})
			return
		}

		reader, writer := io.Pipe()
		go writeSmartTravelReply(ctx, writer, runtime, userID, message)

		c.SetStatusCode(consts.StatusOK)
		c.SetContentType("text/plain; charset=utf-8")
		c.Response.Header.Set("Cache-Control", "no-cache")
		c.Response.Header.Set("X-Accel-Buffering", "no")
		c.SetBodyStream(reader, -1)
	}
}

func writeSmartTravelReply(ctx context.Context, writer *io.PipeWriter, runtime *bootstrap.Runtime, userID int64, message string) {
	defer writer.Close()
	writeSmartTravelReplyWithClient(ctx, writer, runtime, userID, message, travelai.NewClient(runtime.Config.AI))
}

type streamingAIClient interface {
	Available() bool
	StreamChatWithMessages(ctx context.Context, messages []travelai.Message, output io.Writer) error
}

func writeSmartTravelReplyWithClient(ctx context.Context, writer io.Writer, runtime *bootstrap.Runtime, userID int64, message string, client streamingAIClient) {
	if client == nil || !client.Available() {
		writeModelUnavailable(ctx, writer, runtime.Logger)
		return
	}
	history := travelSessions.History(userID)
	referenceContext := retrieveTravelKnowledge(ctx, runtime, message)
	messages := buildAIContextMessages(history, message, referenceContext)

	var replyBuffer bytes.Buffer
	streamWriter := io.MultiWriter(writer, &replyBuffer)
	if err := client.StreamChatWithMessages(ctx, messages, streamWriter); err != nil {
		runtime.Logger.Warn("流式调用旅行智能体模型失败", zap.Error(err))
		if replyBuffer.Len() == 0 {
			writeModelUnavailable(ctx, writer, runtime.Logger)
		}
		return
	}

	replyText := strings.TrimSpace(replyBuffer.String())
	if replyText != "" {
		travelSessions.Append(userID, message, replyText)
	}
}

func retrieveTravelKnowledge(ctx context.Context, runtime *bootstrap.Runtime, message string) string {
	if runtime == nil || runtime.Config == nil || !runtime.Config.RAG.Enabled {
		return ""
	}
	retriever := rag.NewLocalRetriever(rag.DefaultDocuments())
	results, err := retriever.Search(ctx, message, runtime.Config.RAG)
	if err != nil {
		runtime.Logger.Warn("检索旅行知识失败", zap.Error(err))
		return ""
	}
	return rag.FormatContext(results)
}

func buildAIContextMessages(history []conversationMessage, message string, referenceContext ...string) []travelai.Message {
	messages := make([]travelai.Message, 0, len(history)+1)
	for _, item := range history {
		if strings.TrimSpace(item.Content) == "" {
			continue
		}
		messages = append(messages, travelai.Message{Role: item.Role, Content: item.Content})
	}
	currentMessage := message
	if len(referenceContext) > 0 {
		if contextText := strings.TrimSpace(referenceContext[0]); contextText != "" {
			currentMessage = contextText + "\n\n用户问题：" + message
		}
	}
	messages = append(messages, travelai.Message{Role: "user", Content: currentMessage})
	return messages
}

func writeModelUnavailable(ctx context.Context, writer io.Writer, log *zap.Logger) {
	select {
	case <-ctx.Done():
		return
	default:
		if _, err := writer.Write([]byte(modelUnavailableReply)); err != nil {
			log.Warn("旅行智能体不可用提示写入失败", zap.Error(err))
		}
	}
}

func writeTextStream(ctx context.Context, writer io.Writer, text string, log *zap.Logger) {
	for _, piece := range []rune(text) {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := writer.Write([]byte(string(piece))); err != nil {
				log.Warn("旅行智能体回复写入失败", zap.Error(err))
				return
			}
		}
	}
}
