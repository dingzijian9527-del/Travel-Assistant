package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const loginOutcomeKey = "login:outcome"

type loginOutcome string

const (
	loginSucceeded loginOutcome = "succeeded"
	loginFailed    loginOutcome = "failed"
)

func LoginRateLimit() app.HandlerFunc {
	limiter := NewLoginFailureLimiter(maxLoginFailures, loginFailureWindow)
	return func(ctx context.Context, c *app.RequestContext) {
		identity := loginIdentityFromRequest(c)
		if identity.Phone == "" {
			c.Next(ctx)
			return
		}
		if !limiter.Allow(identity) {
			c.AbortWithStatusJSON(consts.StatusTooManyRequests, map[string]any{
				"code": 429,
				"msg":  "登录次数过多，请稍后再试",
			})
			return
		}

		c.Next(ctx)
		outcome, _ := c.Get(loginOutcomeKey)
		switch outcome {
		case loginFailed:
			limiter.RecordFailure(identity)
		case loginSucceeded:
			limiter.Reset(identity)
		}
	}
}

func MarkLoginFailure(c *app.RequestContext) {
	c.Set(loginOutcomeKey, loginFailed)
}

func MarkLoginSuccess(c *app.RequestContext) {
	c.Set(loginOutcomeKey, loginSucceeded)
}

func loginIdentityFromRequest(c *app.RequestContext) loginIdentity {
	return loginIdentity{
		Phone:    normalizePhone(loginPhoneFromRequest(c)),
		IP:       clientIP(c),
		DeviceID: string(c.GetHeader("X-Device-ID")),
	}
}
