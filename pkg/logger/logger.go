package logger

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// New 创建带服务名、环境和本地文件输出的日志器。
func New(service string, cfg config.Config) (*zap.Logger, error) {
	if err := os.MkdirAll(cfg.Log.Dir, 0o755); err != nil {
		return nil, err
	}
	if err := CleanExpired(cfg.Log.Dir, cfg.Log.MaxAgeDays, time.Now()); err != nil {
		return nil, err
	}
	level := zapcore.InfoLevel
	if err := level.Set(cfg.Log.Level); err != nil {
		return nil, err
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	fileName := filepath.Join(cfg.Log.Dir, service+"-"+time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	writers := []zapcore.WriteSyncer{zapcore.AddSync(file)}
	if cfg.Log.Console {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).With(
		zap.String("service", service),
		zap.String("env", cfg.App.Env),
	), nil
}

// CleanExpired 删除超过保留天数的日志文件。
func CleanExpired(dir string, maxAgeDays int, now time.Time) error {
	if maxAgeDays <= 0 {
		maxAgeDays = 7
	}
	deadline := now.Add(-time.Duration(maxAgeDays) * 24 * time.Hour)
	return filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry == nil || entry.IsDir() {
			return err
		}
		if strings.ToLower(filepath.Ext(path)) != ".log" {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.ModTime().Before(deadline) {
			return os.Remove(path)
		}
		return nil
	})
}
