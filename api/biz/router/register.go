package router

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/handler"
	gatewaymw "github.com/dingzijian9527-del/Travel-Assistant/api/biz/middleware"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
)

// Register 按网关路由分组注册接口。公开接口无需鉴权，受保护接口使用 RequireAuth 中间件统一拦截。
func Register(h *server.Hertz, runtime *bootstrap.Runtime) {
	h.Use(gatewaymw.CORS(), gatewaymw.RequestID(), gatewaymw.AccessLog(runtime.Logger), gatewaymw.Recovery(runtime.Logger))
	h.GET("/ping", handler.Ping(runtime))

	v1 := h.Group("/api/v1")
	// 公开接口（无需鉴权）
	v1.GET("/config", handler.ConfigSummary(runtime))
	v1.POST("/sms/register-code", gatewaymw.RegisterCodeRateLimit(), handler.SendRegisterCode(runtime))
	v1.POST("/user/register", gatewaymw.RateLimitByIP("user-register", 5, time.Minute), handler.RegisterUser(runtime))
	v1.POST("/user/login", gatewaymw.RateLimitByIP("user-login", 5, time.Minute), gatewaymw.LoginRateLimit(), handler.LoginUser(runtime))
	v1.POST("/user/refresh", handler.TokenRefresh(runtime))

	// 受保护接口（需要有效令牌）
	auth := v1.Group("", gatewaymw.RequireAuth(runtime))
	{
		auth.GET("/user/profile", handler.GetUserProfile(runtime))
		auth.POST("/user/profile", handler.UpdateUserProfile(runtime))
		auth.GET("/user/dashboard", handler.GetUserDashboard(runtime))
		auth.GET("/user/preferences", handler.GetUserPreferences(runtime))
		auth.POST("/user/preferences", handler.UpdateUserPreferences(runtime))
		auth.GET("/user/settings", handler.GetUserSettings(runtime))
		auth.POST("/user/settings", handler.UpdateUserSettings(runtime))
		auth.POST("/user/confirm", handler.PasswordConfirm(runtime))
		auth.POST("/upload/check", gatewaymw.RateLimitByIP("upload-check", 20, time.Minute), handler.CheckUpload(runtime))
		auth.POST("/upload/avatar", gatewaymw.RateLimitByIP("upload-avatar", 20, time.Minute), handler.UploadAvatar(runtime))
		auth.POST("/ai-stream", gatewaymw.RateLimitByIP("ai-stream", 20, time.Minute), handler.ChatStream(runtime))
		auth.POST("/ai/chat/stream", gatewaymw.RateLimitByIP("ai-chat-stream", 20, time.Minute), handler.ChatStream(runtime))
		auth.POST("/trips", handler.CreateTrip(runtime))
		auth.GET("/trips", handler.ListTrips(runtime))
		auth.GET("/trips/latest", handler.GetLatestTrip(runtime))
		auth.GET("/trips/:tripID", handler.GetTripDetail(runtime))
		auth.DELETE("/trips/:tripID", handler.DeleteTrip(runtime))
		auth.PUT("/trips/:tripID", handler.UpdateTrip(runtime))
	}
}
