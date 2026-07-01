package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	aiagentservice "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent/aiagentservice"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	rpcaiagent "github.com/dingzijian9527-del/Travel-Assistant/rpc/ai_agent"
	"go.uber.org/zap"
)

func main() {
	runtime, err := bootstrap.Init("ai-agent-service")
	if err != nil {
		panic(err)
	}
	defer runtime.Logger.Sync()

	addr := &net.TCPAddr{IP: net.ParseIP(runtime.Config.RPC.Host), Port: runtime.Config.RPC.AIAgent.Port}
	svr := aiagentservice.NewServer(new(rpcaiagent.AIAgentServiceImpl), server.WithServiceAddr(addr))
	runtime.Logger.Info("旅行智能体 Kitex 服务启动", zap.String("addr", addr.String()))
	if err := svr.Run(); err != nil {
		runtime.Logger.Fatal("旅行智能体 Kitex 服务异常退出", zap.Error(err))
	}
}
