package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/redisx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/smsx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/verifycodex"
)

type sendRegisterCodeRequest struct {
	Phone string `json:"phone"`
}

var registerCodeStores sync.Map

// SendRegisterCode 发送注册短信验证码，并把验证码写入缓存服务五分钟。
func SendRegisterCode(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req sendRegisterCodeRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		phone := normalizePhone(req.Phone)
		if !validMainlandPhone(phone) {
			writeJSON(c, consts.StatusBadRequest, 400, "手机号格式错误", nil)
			return
		}
		code, err := generateRegisterCode()
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "验证码生成失败", nil)
			return
		}
		expire := runtime.Config.SMS.RegisterCodeExpire
		if expire <= 0 {
			expire = 5 * time.Minute
		}
		store, err := newRegisterCodeStore(runtime)
		if err != nil {
			runtime.Logger.Warn("注册验证码缓存组件初始化失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, "验证码缓存不可用", nil)
			return
		}
		if err := store.Save(ctx, phone, code, expire); err != nil {
			runtime.Logger.Warn("注册短信验证码缓存失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, "验证码缓存失败", nil)
			return
		}
		if shouldUseLocalRegisterCode(runtime) {
			writeJSON(c, consts.StatusOK, 0, "开发环境验证码已生成", map[string]any{
				"code":           code,
				"expire_seconds": int(expire.Seconds()),
			})
			return
		}
		if !tencentSMSConfigured(runtime.Config.SMS) {
			_ = store.Delete(ctx, phone)
			runtime.Logger.Warn("腾讯云短信配置未完整，无法发送验证码")
			writeJSON(c, consts.StatusInternalServerError, 500, "腾讯云短信配置未完整，无法发送验证码", nil)
			return
		}
		sender := smsx.NewTencentSender(smsx.Config{
			SecretID:   runtime.Config.SMS.SecretID,
			SecretKey:  runtime.Config.SMS.SecretKey,
			SDKAppID:   runtime.Config.SMS.SDKAppID,
			SignName:   runtime.Config.SMS.SignName,
			TemplateID: runtime.Config.SMS.TemplateID,
			Region:     runtime.Config.SMS.Region,
			Endpoint:   runtime.Config.SMS.Endpoint,
		})
		if err := sender.SendRegisterCode(ctx, phone, code, expire); err != nil {
			_ = store.Delete(ctx, phone)
			runtime.Logger.Warn("发送注册短信验证码失败", zap.Error(err))
			writeJSON(c, consts.StatusInternalServerError, 500, "短信验证码发送失败", nil)
			return
		}
		writeJSON(c, consts.StatusOK, 0, "验证码已发送", nil)
	}
}

func newRegisterCodeStore(runtime *bootstrap.Runtime) (verifycodex.Store, error) {
	cacheKey := fmt.Sprintf("%s|%s|%d", runtime.Config.Redis.Addr, runtime.Config.Redis.Username, runtime.Config.Redis.DB)
	if value, ok := registerCodeStores.Load(cacheKey); ok {
		return value.(verifycodex.Store), nil
	}
	client, err := redisx.New(runtime.Config.Redis)
	if err != nil {
		return nil, err
	}
	store := verifycodex.NewRedisStore(client)
	actual, _ := registerCodeStores.LoadOrStore(cacheKey, store)
	return actual.(verifycodex.Store), nil
}

func shouldUseLocalRegisterCode(runtime *bootstrap.Runtime) bool {
	if runtime == nil || runtime.Config == nil {
		return false
	}
	return strings.EqualFold(runtime.Config.App.Env, "dev") &&
		runtime.Config.SMS.DevReturnCode &&
		!tencentSMSConfigured(runtime.Config.SMS)
}

func tencentSMSConfigured(cfg config.SMSConfig) bool {
	return strings.TrimSpace(cfg.SecretID) != "" &&
		strings.TrimSpace(cfg.SecretKey) != "" &&
		strings.TrimSpace(cfg.SDKAppID) != "" &&
		strings.TrimSpace(cfg.SignName) != "" &&
		strings.TrimSpace(cfg.TemplateID) != ""
}

func generateRegisterCode() (string, error) {
	max := big.NewInt(1000000)
	value, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", value.Int64()), nil
}

func normalizePhone(phone string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(phone, "")
}

func validMainlandPhone(phone string) bool {
	return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(phone)
}
