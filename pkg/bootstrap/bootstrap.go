package bootstrap

import (
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/logger"
)

var defaultConfigPath string

func init() {
	defaultConfigPath = strings.TrimSpace(os.Getenv("TRAVEL_ASSISTANT_CONFIG"))
	if defaultConfigPath == "" {
		defaultConfigPath = "conf/config.yaml"
	}
}

// Runtime 保存进程启动时预加载的配置和日志组件。
type Runtime struct {
	Config *config.Config
	Logger *zap.Logger
}

// Init 在 main 中统一初始化所有需要提前加载的基础组件。
func Init(service string) (*Runtime, error) {
	return InitWithConfigPath(service, defaultConfigPath)
}

// InitWithConfigPath 允许测试或脚本显式指定配置文件路径。
func InitWithConfigPath(service string, path string) (*Runtime, error) {
	cfg, err := config.InitGlobal(path)
	if err != nil {
		return nil, err
	}
	if err := cfg.ValidateForService(service); err != nil {
		return nil, err
	}
	log, err := logger.New(service, *cfg)
	if err != nil {
		return nil, err
	}
	log.Info("基础组件初始化完成", zap.String("service", service))
	return &Runtime{Config: cfg, Logger: log}, nil
}
