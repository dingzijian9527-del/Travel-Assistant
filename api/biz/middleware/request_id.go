package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID 为每个网关请求注入请求编号。
func RequestID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := string(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Response.Header.Set(requestIDHeader, requestID)
		c.Next(ctx)
	}
}
