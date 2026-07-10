package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	apiupload "github.com/dingzijian9527-del/Travel-Assistant/api/biz/utils/upload"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/uploadx"
)

func UploadAvatar(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if _, ok := claimsFromRequest(runtime, c); !ok {
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil || fileHeader == nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请选择头像图片", nil)
			return
		}
		if err := apiupload.Check(fileHeader.Filename, fileHeader.Size, runtime.Config.Upload); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, err.Error(), nil)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "读取头像文件失败", nil)
			return
		}
		defer file.Close()

		url, key, err := uploadx.UploadAvatarToQiniu(ctx, runtime.Config.Upload.Qiniu, fileHeader.Filename, file)
		if err != nil {
			runtime.Logger.Warn("上传头像到七牛云失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, friendlyUploadError(err), nil)
			return
		}

		writeJSON(c, consts.StatusOK, 0, "success", map[string]any{
			"url": url,
			"key": key,
		})
	}
}

func friendlyUploadError(err error) string {
	if err == nil {
		return "上传失败"
	}
	message := strings.TrimSpace(err.Error())
	if message == "" {
		return "上传失败"
	}
	if strings.Contains(message, "未配置") || strings.Contains(message, "文件") || strings.Contains(message, "不支持") {
		return message
	}
	return fmt.Sprintf("头像上传失败：%s", message)
}
