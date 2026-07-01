package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

// Recovery 统一拦截网关异常，避免错误穿透到调用方。
func Recovery(log *zap.Logger) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("HTTP 请求发生异常", zap.Any("panic", err))
				c.AbortWithStatusJSON(consts.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统内部错误"})
			}
		}()
		c.Next(ctx)
	}
}
