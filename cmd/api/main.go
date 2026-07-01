package main

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	apiapp "github.com/dingzijian9527-del/Travel-Assistant/api"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"go.uber.org/zap"
)

func main() {
	runtime, err := bootstrap.Init("api-gateway")
	if err != nil {
		panic(err)
	}
	defer runtime.Logger.Sync()

	cfg := runtime.Config
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	h := server.Default(
		server.WithHostPorts(addr),
		server.WithReadTimeout(cfg.HTTP.ReadTimeout),
		server.WithWriteTimeout(cfg.HTTP.WriteTimeout),
	)
	apiapp.Register(h, runtime)
	runtime.Logger.Info("Hertz 网关启动", zap.String("addr", addr))
	h.Spin()
}
