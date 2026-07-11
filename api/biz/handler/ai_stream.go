package handler

import (
	"context"
	"io"
	"unicode/utf8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/middleware"
	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/validator"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
	rpcaiagent "github.com/dingzijian9527-del/Travel-Assistant/rpc/ai_agent"
)

const modelUnavailableReply = "智能体模型暂不可用，请稍后重试或检查火山方舟接入点配置。"

// ChatStreamRequest 描述旅行智能体流式对话请求。
type ChatStreamRequest struct {
	Message string `json:"message"`
}

// ChatStream 处理旅行智能体流式对话。
// 通过 HTTP 流式输出智能体回复，内部调用 rpc/ai_agent 包的 ChatStream 方法。
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

		claims, ok := middleware.AuthClaims(c)
		if !ok {
			return
		}
		userID, err := jwtx.UserIDInt64(claims.UserID)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]any{"code": 401, "msg": "登录状态无效，请重新登录"})
			return
		}

		reader, writer := io.Pipe()
		go writeChatStreamReply(ctx, writer, runtime, userID, message)

		c.SetStatusCode(consts.StatusOK)
		c.SetContentType("text/plain; charset=utf-8")
		c.Response.Header.Set("Cache-Control", "no-cache")
		c.Response.Header.Set("X-Accel-Buffering", "no")
		c.SetBodyStream(reader, -1)
	}
}

func writeChatStreamReply(ctx context.Context, writer *io.PipeWriter, runtime *bootstrap.Runtime, userID int64, message string) {
	defer writer.Close()

	aiAgent := newAIAgentStreamService(runtime)
	if aiAgent == nil {
		_, _ = writer.Write([]byte(unavailableServiceReply()))
		return
	}
	svcErr := aiAgent.ChatStream(ctx, userID, message, writer)
	if svcErr != nil {
		runtime.Logger.Warn("旅行智能体流式对话失败", zap.Error(svcErr))
	}
}

func newAIAgentStreamService(runtime *bootstrap.Runtime) *rpcaiagent.AIAgentStreamService {
	service, err := rpcaiagent.NewAIAgentStreamService(
		runtime.Config.RAG,
		runtime.Config.AI,
		runtime.Config.MySQL,
		runtime.Config.TravelData,
		runtime.Logger,
	)
	if err != nil {
		runtime.Logger.Warn("旅行智能体流式服务初始化失败", zap.Error(err))
		return nil
	}
	return service
}

func unavailableServiceReply() string {
	return "智能体服务初始化失败，请稍后重试。"
}
