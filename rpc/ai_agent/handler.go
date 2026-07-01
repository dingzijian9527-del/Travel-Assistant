package rpcaiagent

import (
	"context"
	"sync"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// AIAgentServiceImpl 实现接口定义中的旅行智能体服务。
type AIAgentServiceImpl struct {
	initOnce sync.Once
	service  *aiAgentService
}

func (s *AIAgentServiceImpl) getService() *aiAgentService {
	s.initOnce.Do(func() {
		s.service = newAIAgentServiceWithModel(newAIAgentRepo(), newTravelModelClient(config.MustGlobal().AI))
	})
	return s.service
}

// Chat 实现旅行智能体对话接口。
func (s *AIAgentServiceImpl) Chat(ctx context.Context, req *aiagent.ChatRequest) (*aiagent.ChatResponse, error) {
	if req == nil {
		return chatErrorResponse(errParam("request is required")), nil
	}
	reply, suggestions, svcErr := s.getService().Chat(ctx, req)
	if svcErr != nil {
		return chatErrorResponse(svcErr), nil
	}
	return &aiagent.ChatResponse{BaseResp: successBaseResp(), Reply: toChatMessageDTO(reply), Suggestions: suggestions}, nil
}

// GetPromptSuggestions 实现旅行提示词推荐接口。
func (s *AIAgentServiceImpl) GetPromptSuggestions(ctx context.Context, req *aiagent.PromptSuggestionsRequest) (*aiagent.PromptSuggestionsResponse, error) {
	if req == nil {
		req = aiagent.NewPromptSuggestionsRequest()
	}
	suggestions := s.getService().GetPromptSuggestions(ctx, req)
	return &aiagent.PromptSuggestionsResponse{BaseResp: successBaseResp(), Suggestions: suggestions}, nil
}
