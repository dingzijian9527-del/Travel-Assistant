package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// AccessLog 统一打印赫兹访问日志。
func AccessLog(log *zap.Logger) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		log.Info("HTTP 请求完成",
			zap.String("request_id", requestID(c)),
			zap.String("method", string(c.Method())),
			zap.String("path", string(c.Path())),
			zap.Int("status", c.Response.StatusCode()),
			zap.Duration("cost", time.Since(start)),
		)
	}
}

func requestID(c *app.RequestContext) string {
	value, ok := c.Get("request_id")
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}
