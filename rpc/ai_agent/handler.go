package rpcaiagent

import (
	"context"
	"io"
	"sync"
	"time"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/traveldata"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

// AIAgentServiceImpl 实现 AI 智能体 RPC 服务。
type AIAgentServiceImpl struct {
	initOnce sync.Once
	service  *aiAgentService
	initErr  error
	logger   *zap.Logger
}

// NewAIAgentServiceImpl 创建 AI 智能体服务实现。
func NewAIAgentServiceImpl() *AIAgentServiceImpl {
	return &AIAgentServiceImpl{}
}

// getService 初始化 AI 智能体服务（延迟初始化），包括连接 Milvus 或加载本地检索器。
func (s *AIAgentServiceImpl) getService(ctx context.Context) (*aiAgentService, error) {
	s.initOnce.Do(func() {
		cfg := config.MustGlobal()
		// 初始化日志
		logger, _ := zap.NewProduction()
		s.logger = logger

		// 创建检索器
		retriever, err := s.buildRetriever(ctx, cfg)
		if err != nil {
			// 检索器初始化失败不阻塞服务启动，用空检索器（不返回任何知识）
			if logger != nil {
				logger.Warn("AI 智能体检索器初始化失败，将使用空检索器", zap.Error(err))
			}
			retriever = emptyRetriever{}
		}

		repo, err := mysqlrepo.NewAIAgentMessageRepository(cfg.MySQL)
		if err != nil {
			s.initErr = err
			return
		}

		s.service = newAIAgentServiceWithTravelData(
			repo,
			newTravelModelClient(cfg.AI),
			retriever,
			cfg.RAG,
			traveldata.NewPlannerFromConfig(cfg.TravelData),
		)
	})
	return s.service, s.initErr
}

// buildRetriever 根据配置创建检索器。优先连接 Milvus，失败时使用本地关键词检索。
func (s *AIAgentServiceImpl) buildRetriever(ctx context.Context, cfg *config.Config) (rag.Retriever, error) {
	if !cfg.RAG.Enabled {
		return nil, nil
	}

	// 优先创建 Milvus 向量检索器
	if cfg.RAG.Address != "" {
		connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		milvusClient, err := client.NewClient(connectCtx, client.Config{Address: cfg.RAG.Address})
		if err == nil {
			embedder, embedderErr := rag.NewConfiguredEmbedder(cfg.RAG)
			if embedderErr != nil {
				return nil, embedderErr
			}
			if s.logger != nil {
				s.logger.Info("AI 智能体已连接 Milvus 向量数据库", zap.String("address", cfg.RAG.Address))
			}
			return rag.NewMilvusRetriever(milvusClient, embedder), nil
		}
		if s.logger != nil {
			s.logger.Warn("无法连接 Milvus，降级为本地关键词检索", zap.String("address", cfg.RAG.Address), zap.Error(err))
		}
	}

	// 降级：使用本地关键词检索
	if s.logger != nil {
		s.logger.Info("AI 智能体使用本地关键词检索")
	}
	return rag.NewLocalRetriever(rag.DefaultDocuments()), nil
}

// Chat 实现 AI 智能体对话 RPC 接口。
func (s *AIAgentServiceImpl) Chat(ctx context.Context, req *aiagent.ChatRequest) (*aiagent.ChatResponse, error) {
	if req == nil {
		return chatErrorResponse(errParam("request is required")), nil
	}
	service, err := s.getService(ctx)
	if err != nil {
		return chatErrorResponse(errParam("AI 智能体服务初始化失败")), nil
	}
	reply, suggestions, svcErr := service.Chat(ctx, req)
	if svcErr != nil {
		return chatErrorResponse(svcErr), nil
	}
	return &aiagent.ChatResponse{BaseResp: successBaseResp(), Reply: toChatMessageDTO(reply), Suggestions: suggestions}, nil
}

func (s *AIAgentServiceImpl) ChatStream(ctx context.Context, req *aiagent.ChatRequest, stream aiagent.AIAgentService_ChatStreamServer) error {
	if req == nil {
		return stream.Send(ctx, &aiagent.ChatStreamChunk{BaseResp: errorBaseResp(errParam("request is required")), Done: true})
	}
	service, err := s.getService(ctx)
	if err != nil {
		return stream.Send(ctx, &aiagent.ChatStreamChunk{BaseResp: errorBaseResp(errParam("AI 智能体服务初始化失败")), Done: true})
	}
	writer := &chatStreamChunkWriter{ctx: ctx, stream: stream}
	if svcErr := service.ChatStream(ctx, req.UserId, req.Message, writer); svcErr != nil {
		return stream.Send(ctx, &aiagent.ChatStreamChunk{BaseResp: errorBaseResp(svcErr), Done: true})
	}
	return stream.Send(ctx, &aiagent.ChatStreamChunk{BaseResp: successBaseResp(), Done: true})
}

type chatStreamChunkWriter struct {
	ctx    context.Context
	stream aiagent.AIAgentService_ChatStreamServer
}

func (w *chatStreamChunkWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	content := string(data)
	err := w.stream.Send(w.ctx, &aiagent.ChatStreamChunk{BaseResp: successBaseResp(), Content: &content, Done: false})
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

var _ io.Writer = (*chatStreamChunkWriter)(nil)

// GetPromptSuggestions 实现提示词推荐 RPC 接口。
func (s *AIAgentServiceImpl) GetPromptSuggestions(ctx context.Context, req *aiagent.PromptSuggestionsRequest) (*aiagent.PromptSuggestionsResponse, error) {
	if req == nil {
		req = aiagent.NewPromptSuggestionsRequest()
	}
	service, err := s.getService(ctx)
	if err != nil {
		return &aiagent.PromptSuggestionsResponse{BaseResp: errorBaseResp(errParam("AI 智能体服务初始化失败"))}, nil
	}
	suggestions := service.GetPromptSuggestions(ctx, req)
	return &aiagent.PromptSuggestionsResponse{BaseResp: successBaseResp(), Suggestions: suggestions}, nil
}

// emptyRetriever 是一个不返回任何结果的空检索器，用于检索器初始化失败时的兜底。
type emptyRetriever struct{}

func (e emptyRetriever) Search(ctx context.Context, query string, cfg config.RAGConfig) ([]rag.Result, error) {
	return nil, nil
}
