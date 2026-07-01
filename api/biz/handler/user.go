package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

type registerRequest struct {
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	HomeCity    string `json:"home_city"`
	CurrentCity string `json:"current_city"`
}

type loginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type updateProfileRequest struct {
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	HomeCity    string `json:"home_city"`
	CurrentCity string `json:"current_city"`
}

func RegisterUser(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req registerRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.Register(ctx, &user.RegisterRequest{
			Phone:       req.Phone,
			Password:    req.Password,
			Nickname:    optionalRequestString(req.Nickname),
			AvatarUrl:   optionalRequestString(req.AvatarURL),
			HomeCity:    optionalRequestString(req.HomeCity),
			CurrentCity: optionalRequestString(req.CurrentCity),
		})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetUser())
	}
}

func LoginUser(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req loginRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.Login(ctx, &user.LoginRequest{Phone: req.Phone, Password: req.Password})
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		if resp.GetBaseResp().GetCode() != 0 {
			writeJSON(c, consts.StatusOK, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), nil)
			return
		}
		writeJSON(c, consts.StatusOK, 0, "success", map[string]any{"token": resp.GetToken(), "user": resp.GetUser()})
	}
}

func GetUserProfile(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, ok := claimsFromRequest(runtime, c)
		if !ok {
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.GetProfile(ctx, &user.GetProfileRequest{Id: claims.UserID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetUser())
	}
}

func UpdateUserProfile(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, ok := claimsFromRequest(runtime, c)
		if !ok {
			return
		}
		var req updateProfileRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.UpdateProfile(ctx, &user.UpdateProfileRequest{
			Id:          claims.UserID,
			Nickname:    optionalRequestString(req.Nickname),
			AvatarUrl:   optionalRequestString(req.AvatarURL),
			HomeCity:    optionalRequestString(req.HomeCity),
			CurrentCity: optionalRequestString(req.CurrentCity),
		})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetUser())
	}
}

func claimsFromRequest(runtime *bootstrap.Runtime, c *app.RequestContext) (jwtx.Claims, bool) {
	token := bearerToken(string(c.Request.Header.Peek("Authorization")))
	claims, err := jwtx.Parse(jwtx.Config{Secret: runtime.Config.Auth.JWTSecret, Expire: runtime.Config.Auth.JWTExpire}, token)
	if err != nil {
		writeJSON(c, consts.StatusUnauthorized, 401, "登录状态无效，请重新登录", nil)
		return jwtx.Claims{}, false
	}
	return claims, true
}

func bearerToken(authorization string) string {
	parts := strings.Fields(authorization)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func optionalRequestString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func writeRPCResponse(c *app.RequestContext, code int32, msg string, data any) {
	writeJSON(c, consts.StatusOK, code, msg, data)
}

func writeJSON(c *app.RequestContext, status int, code int32, msg string, data any) {
	body := map[string]any{"code": code, "msg": msg}
	if data != nil {
		body["data"] = data
	}
	c.JSON(status, body)
}
