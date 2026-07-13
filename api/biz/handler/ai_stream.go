package handler

import (
	"context"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/middleware"
	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/validator"
	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

const modelUnavailableReply = "智能体模型暂不可用，请稍后重试或检查火山方舟接入点配置。"

type ChatStreamRequest struct {
	Message string `json:"message"`
}

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

	clients, err := clientsFor(runtime)
	if err != nil {
		_, _ = writer.Write([]byte(unavailableServiceReply()))
		return
	}
	stream, err := clients.aiAgent.ChatStream(ctx, &aiagent.ChatRequest{UserId: userID, Message: message})
	if err != nil {
		_, _ = writer.Write([]byte(unavailableServiceReply()))
		return
	}
	for {
		chunk, err := stream.Recv(ctx)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil || chunk == nil || chunk.GetBaseResp().GetCode() != 0 {
			_, _ = writer.Write([]byte(unavailableServiceReply()))
			return
		}
		if content := chunk.GetContent(); content != "" {
			if _, err := writer.Write([]byte(content)); err != nil {
				return
			}
		}
		if chunk.GetDone() {
			return
		}
	}
}

func unavailableServiceReply() string {
	return "智能体服务暂不可用，请稍后重试。"
}
