package router

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/handler"
	gatewaymw "github.com/dingzijian9527-del/Travel-Assistant/api/biz/middleware"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
)

// Register 按网关路由分组注册接口。
func Register(h *server.Hertz, runtime *bootstrap.Runtime) {
	h.Use(gatewaymw.CORS(), gatewaymw.RequestID(), gatewaymw.AccessLog(runtime.Logger), gatewaymw.Recovery(runtime.Logger))
	h.GET("/ping", handler.Ping(runtime))
	v1 := h.Group("/api/v1")
	{
		v1.GET("/config", handler.ConfigSummary(runtime))
		v1.POST("/sms/register-code", gatewaymw.RegisterCodeRateLimit(), handler.SendRegisterCode(runtime))
		v1.POST("/user/register", gatewaymw.RateLimitByIP("user-register", 5, time.Minute), handler.RegisterUser(runtime))
		v1.POST("/user/login", gatewaymw.RateLimitByIP("user-login", 5, time.Minute), handler.LoginUser(runtime))
		v1.GET("/user/profile", handler.GetUserProfile(runtime))
		v1.POST("/user/profile", handler.UpdateUserProfile(runtime))
		v1.GET("/user/dashboard", handler.GetUserDashboard(runtime))
		v1.GET("/user/preferences", handler.GetUserPreferences(runtime))
		v1.POST("/user/preferences", handler.UpdateUserPreferences(runtime))
		v1.GET("/user/settings", handler.GetUserSettings(runtime))
		v1.POST("/user/settings", handler.UpdateUserSettings(runtime))
		v1.POST("/upload/check", gatewaymw.RateLimitByIP("upload-check", 20, time.Minute), handler.CheckUpload(runtime))
		v1.POST("/upload/avatar", gatewaymw.RateLimitByIP("upload-avatar", 20, time.Minute), handler.UploadAvatar(runtime))
		v1.POST("/ai-stream", gatewaymw.RateLimitByIP("ai-stream", 20, time.Minute), handler.ChatStream(runtime))
		v1.POST("/ai/chat/stream", gatewaymw.RateLimitByIP("ai-chat-stream", 20, time.Minute), handler.ChatStream(runtime))
		v1.POST("/trips", handler.CreateTrip(runtime))
		v1.GET("/trips", handler.ListTrips(runtime))
		v1.GET("/trips/latest", handler.GetLatestTrip(runtime))
		v1.GET("/trips/:tripID", handler.GetTripDetail(runtime))
		v1.DELETE("/trips/:tripID", handler.DeleteTrip(runtime))
	}
}
