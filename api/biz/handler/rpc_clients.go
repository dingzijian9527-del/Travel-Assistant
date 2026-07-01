package handler

import (
	"errors"
	"sync"

	kitexclient "github.com/cloudwego/kitex/client"
	aiagentservice "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent/aiagentservice"
	userservice "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user/userservice"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
)

type rpcClients struct {
	user    userservice.Client
	aiAgent aiagentservice.Client
}

var gatewayClients sync.Map

func clientsFor(runtime *bootstrap.Runtime) (*rpcClients, error) {
	if runtime == nil || runtime.Config == nil {
		return nil, errors.New("runtime config is required")
	}
	cacheKey := runtime.Config.RPC.User.Target + "|" + runtime.Config.RPC.AIAgent.Target
	if value, ok := gatewayClients.Load(cacheKey); ok {
		return value.(*rpcClients), nil
	}
	userClient, err := userservice.NewClient(runtime.Config.RPC.User.ServiceName, kitexclient.WithHostPorts(runtime.Config.RPC.User.Target))
	if err != nil {
		return nil, err
	}
	aiClient, err := aiagentservice.NewClient(runtime.Config.RPC.AIAgent.ServiceName, kitexclient.WithHostPorts(runtime.Config.RPC.AIAgent.Target))
	if err != nil {
		return nil, err
	}
	clients := &rpcClients{user: userClient, aiAgent: aiClient}
	actual, _ := gatewayClients.LoadOrStore(cacheKey, clients)
	return actual.(*rpcClients), nil
}
