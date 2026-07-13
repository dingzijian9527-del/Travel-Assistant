package middleware

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	commonconfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

type testHeader struct {
	Key   string
	Value string
}

func TestFixedWindowLimiterBlocksAfterLimitAndResets(t *testing.T) {
	limiter := NewFixedWindowLimiter(2, time.Minute)
	base := time.Date(2026, 7, 10, 12, 0, 0, 0, time.Local)
	limiter.now = func() time.Time { return base }

	if !limiter.Allow("ip:127.0.0.1") {
		t.Fatal("first request should pass")
	}
	if !limiter.Allow("ip:127.0.0.1") {
		t.Fatal("second request should pass")
	}
	if limiter.Allow("ip:127.0.0.1") {
		t.Fatal("third request in same window should be blocked")
	}

	base = base.Add(time.Minute + time.Second)
	if !limiter.Allow("ip:127.0.0.1") {
		t.Fatal("request after window reset should pass")
	}
}

func TestRateLimitByIPReturnsTooManyRequests(t *testing.T) {
	engine := route.NewEngine(commonconfig.NewOptions(nil))
	engine.POST("/login", RateLimitByIP("login", 1, time.Minute), func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "ok")
	})

	first := performRequest(engine, consts.MethodPost, "/login", nil, testHeader{Key: "X-Forwarded-For", Value: "127.0.0.1"})
	if first.Code != consts.StatusOK {
		t.Fatalf("first request status = %d, want %d", first.Code, consts.StatusOK)
	}

	second := performRequest(engine, consts.MethodPost, "/login", nil, testHeader{Key: "X-Forwarded-For", Value: "127.0.0.1"})
	if second.Code != consts.StatusTooManyRequests {
		t.Fatalf("second request status = %d, want %d", second.Code, consts.StatusTooManyRequests)
	}
}

func TestRegisterCodeRateLimitBlocksRepeatedPhoneRequests(t *testing.T) {
	engine := route.NewEngine(commonconfig.NewOptions(nil))
	hitCount := 0
	engine.POST("/sms/register-code", RegisterCodeRateLimit(), func(ctx context.Context, c *app.RequestContext) {
		hitCount++
		c.String(consts.StatusOK, "ok")
	})

	body := []byte(`{"phone":"13800138000"}`)
	headers := []testHeader{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "X-Forwarded-For", Value: "127.0.0.1"},
		{Key: "X-Device-ID", Value: "device-1"},
	}
	first := performRequest(engine, consts.MethodPost, "/sms/register-code", body, headers...)
	if first.Code != consts.StatusOK {
		t.Fatalf("first request status = %d, want %d", first.Code, consts.StatusOK)
	}
	second := performRequest(engine, consts.MethodPost, "/sms/register-code", body, headers...)
	if second.Code != consts.StatusTooManyRequests {
		t.Fatalf("second request status = %d, want %d", second.Code, consts.StatusTooManyRequests)
	}
	if hitCount != 1 {
		t.Fatalf("handler should only run once, got %d", hitCount)
	}
}

func performRequest(engine *route.Engine, method string, path string, body []byte, headers ...testHeader) *httptest.ResponseRecorder {
	ctx := engine.NewContext()
	request := protocol.NewRequest(method, path, nil)
	request.CopyTo(&ctx.Request)
	if body != nil {
		ctx.Request.SetBody(body)
	}
	for _, item := range headers {
		ctx.Request.Header.Set(item.Key, item.Value)
	}

	engine.ServeHTTP(context.Background(), ctx)

	writer := httptest.NewRecorder()
	ctx.Response.Header.VisitAll(func(key, value []byte) {
		writer.Header().Add(string(key), string(value))
	})
	writer.WriteHeader(ctx.Response.StatusCode())
	_, _ = writer.Write(ctx.Response.Body())
	ctx.Reset()
	return writer
}

