package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

const authClaimsKey = "auth:claims"

// AuthClaims 鉴权中间件解析出的用户身份。
func AuthClaims(c *app.RequestContext) (jwtx.Claims, bool) {
	value, ok := c.Get(authClaimsKey)
	if !ok {
		return jwtx.Claims{}, false
	}
	claims, ok := value.(jwtx.Claims)
	return claims, ok
}

// RequireAuth 要求请求携带有效令牌，否则返回 401。
func RequireAuth(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, ok := parseAuthToken(runtime.Config.Auth, c)
		if !ok {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]any{
				"code": 401,
				"msg":  "登录状态无效，请重新登录",
			})
			return
		}
		c.Set(authClaimsKey, claims)
		c.Next(ctx)
	}
}

func parseAuthToken(auth config.AuthConfig, c *app.RequestContext) (jwtx.Claims, bool) {
	token := extractBearerToken(string(c.Request.Header.Peek("Authorization")))
	if token == "" {
		return jwtx.Claims{}, false
	}
	claims, err := jwtx.Parse(jwtx.Config{Secret: auth.JWTSecret, Expire: auth.JWTExpire}, token)
	if err != nil {
		return jwtx.Claims{}, false
	}
	return claims, true
}

func extractBearerToken(authorization string) string {
	parts := strings.Fields(authorization)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}
