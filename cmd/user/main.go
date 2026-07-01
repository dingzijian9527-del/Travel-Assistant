package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	userservice "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user/userservice"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	rpcuser "github.com/dingzijian9527-del/Travel-Assistant/rpc/user"
	"go.uber.org/zap"
)

func main() {
	runtime, err := bootstrap.Init("user-service")
	if err != nil {
		panic(err)
	}
	defer runtime.Logger.Sync()

	addr := &net.TCPAddr{IP: net.ParseIP(runtime.Config.RPC.Host), Port: runtime.Config.RPC.User.Port}
	svr := userservice.NewServer(new(rpcuser.UserServiceImpl), server.WithServiceAddr(addr))
	runtime.Logger.Info("用户 Kitex 服务启动", zap.String("addr", addr.String()))
	if err := svr.Run(); err != nil {
		runtime.Logger.Fatal("用户 Kitex 服务异常退出", zap.Error(err))
	}
}
