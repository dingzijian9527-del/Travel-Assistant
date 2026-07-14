package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/utils/upload"
	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/validator"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// UploadCheckRequest 是网关上传校验请求。
type UploadCheckRequest struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// Ping 返回网关健康状态。
func Ping(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]any{"code": 0, "msg": "pong", "service": runtime.Config.App.Name})
	}
}

// ConfigSummary 返回脱敏后的配置摘要。
func ConfigSummary(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]any{
			"code": 0,
			"msg":  "ok",
			"data": config.SafeConfigSummary(runtime.Config),
		})
	}
}

// CheckUpload 演示文件上传工具类的统一调用入口。
func CheckUpload(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req UploadCheckRequest
		if err := c.BindAndValidate(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": "请求参数错误"})
			return
		}
		filename, ok, msg := validator.RequiredString(req.Filename, "filename")
		if !ok {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": msg})
			return
		}
		if err := upload.Check(filename, req.Size, runtime.Config.Upload); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]any{"code": 400, "msg": err.Error()})
			return
		}
		c.JSON(consts.StatusOK, map[string]any{"code": 0, "msg": "文件校验通过"})
	}
}
