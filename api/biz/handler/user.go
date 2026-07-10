package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/verifycodex"
	"go.uber.org/zap"
)

type registerRequest struct {
	Phone       string `json:"phone"`
	Code        string `json:"code"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	HomeCity    string `json:"home_city"`
	CurrentCity string `json:"current_city"`
}

var (
	errRegisterCodeRequired = errors.New("请输入验证码")
	errRegisterCodeInvalid  = errors.New("验证码错误或已过期")
	errRegisterCodeStorage  = errors.New("验证码校验失败")
)

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

type updatePreferencesRequest struct {
	Items []string `json:"items"`
}

type updateSettingsRequest struct {
	TripReminderEnabled          bool `json:"trip_reminder_enabled"`
	PriceReminderEnabled         bool `json:"price_reminder_enabled"`
	PersonalizedRecommendEnabled bool `json:"personalized_recommend_enabled"`
}

func RegisterUser(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req registerRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		phone := normalizePhone(req.Phone)
		if !validMainlandPhone(phone) {
			writeJSON(c, consts.StatusBadRequest, 400, "手机号格式错误", nil)
			return
		}
		store, err := newRegisterCodeStore(runtime)
		if err != nil {
			runtime.Logger.Warn("注册验证码缓存组件初始化失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, "验证码服务暂不可用", nil)
			return
		}
		if err := validateRegisterCode(ctx, store, phone, req.Code); err != nil {
			if errors.Is(err, errRegisterCodeRequired) || errors.Is(err, errRegisterCodeInvalid) {
				writeJSON(c, consts.StatusBadRequest, 400, err.Error(), nil)
				return
			}
			runtime.Logger.Warn("注册验证码校验失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, "验证码校验失败", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.Register(ctx, &user.RegisterRequest{
			Phone:       phone,
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
		if resp.GetBaseResp().GetCode() == 0 {
			if err := store.Delete(ctx, phone); err != nil {
				runtime.Logger.Warn("注册成功后删除验证码失败", zap.Error(err))
			}
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetUser())
	}
}

func validateRegisterCode(ctx context.Context, store verifycodex.Store, phone string, code string) error {
	if strings.TrimSpace(code) == "" {
		return errRegisterCodeRequired
	}
	ok, err := store.Check(ctx, phone, code)
	if err != nil {
		return fmt.Errorf("%w: %v", errRegisterCodeStorage, err)
	}
	if !ok {
		return errRegisterCodeInvalid
	}
	return nil
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

func GetUserDashboard(runtime *bootstrap.Runtime) app.HandlerFunc {
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
		resp, err := clients.user.GetDashboard(ctx, &user.GetDashboardRequest{Id: claims.UserID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetDashboard())
	}
}

func GetUserPreferences(runtime *bootstrap.Runtime) app.HandlerFunc {
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
		resp, err := clients.user.GetPreferences(ctx, &user.GetPreferencesRequest{Id: claims.UserID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetPreferences())
	}
}

func UpdateUserPreferences(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, ok := claimsFromRequest(runtime, c)
		if !ok {
			return
		}
		var req updatePreferencesRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.UpdatePreferences(ctx, &user.UpdatePreferencesRequest{
			Id:    claims.UserID,
			Items: req.Items,
		})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetPreferences())
	}
}

func GetUserSettings(runtime *bootstrap.Runtime) app.HandlerFunc {
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
		resp, err := clients.user.GetSettings(ctx, &user.GetSettingsRequest{Id: claims.UserID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetSettings())
	}
}

func UpdateUserSettings(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, ok := claimsFromRequest(runtime, c)
		if !ok {
			return
		}
		var req updateSettingsRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务暂不可用", nil)
			return
		}
		resp, err := clients.user.UpdateSettings(ctx, &user.UpdateSettingsRequest{
			Id:                           claims.UserID,
			TripReminderEnabled:          req.TripReminderEnabled,
			PriceReminderEnabled:         req.PriceReminderEnabled,
			PersonalizedRecommendEnabled: req.PersonalizedRecommendEnabled,
		})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "用户服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetSettings())
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
