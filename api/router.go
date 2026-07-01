package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/router"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
)

// Register 注册赫兹网关路由。
func Register(h *server.Hertz, runtime *bootstrap.Runtime) {
	router.Register(h, runtime)
}
