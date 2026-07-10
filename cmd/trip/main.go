package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	tripservice "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip/tripservice"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	rpctrip "github.com/dingzijian9527-del/Travel-Assistant/rpc/trip"
	"go.uber.org/zap"
)

func main() {
	runtime, err := bootstrap.Init("trip-service")
	if err != nil {
		panic(err)
	}
	defer runtime.Logger.Sync()

	addr := &net.TCPAddr{IP: net.ParseIP(runtime.Config.RPC.Host), Port: runtime.Config.RPC.Trip.Port}
	svr := tripservice.NewServer(new(rpctrip.TripServiceImpl), server.WithServiceAddr(addr))
	runtime.Logger.Info("行程 Kitex 服务启动", zap.String("addr", addr.String()))
	if err := svr.Run(); err != nil {
		runtime.Logger.Fatal("行程 Kitex 服务异常退出", zap.Error(err))
	}
}
