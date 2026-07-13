package rpcaiagent

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/traveldata"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

type AIAgentStreamService struct {
	service *aiAgentService
}

var (
	sharedStreamRepoOnce sync.Once
	sharedStreamRepo     aiAgentRepository
	sharedStreamRepoErr  error
)

func NewAIAgentStreamService(ragCfg config.RAGConfig, aiCfg config.AIConfig, mysqlCfg config.MySQLConfig, travelCfg config.TravelDataConfig, logger *zap.Logger) (*AIAgentStreamService, error) {
	repo, err := getSharedStreamRepo(mysqlCfg)
	if err != nil {
		return nil, err
	}
	retriever := buildRetriever(ragCfg, logger)
	modelClient := newTravelModelClient(aiCfg)
	svc := &AIAgentStreamService{
		service: newAIAgentServiceWithTravelData(
			repo,
			modelClient,
			retriever,
			ragCfg,
			traveldata.NewPlannerFromConfig(travelCfg),
		),
	}
	return svc, nil
}

func (s *AIAgentStreamService) ChatStream(ctx context.Context, userID int64, message string, replyWriter io.Writer) *serviceError {
	return s.service.ChatStream(ctx, userID, message, replyWriter)
}

func getSharedStreamRepo(mysqlCfg config.MySQLConfig) (aiAgentRepository, error) {
	sharedStreamRepoOnce.Do(func() {
		sharedStreamRepo, sharedStreamRepoErr = mysqlrepo.NewAIAgentMessageRepository(mysqlCfg)
	})
	return sharedStreamRepo, sharedStreamRepoErr
}

func buildRetriever(ragCfg config.RAGConfig, logger *zap.Logger) rag.Retriever {
	if !ragCfg.Enabled {
		return nil
	}

	if ragCfg.Address != "" {
		connectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		milvusClient, err := client.NewClient(connectCtx, client.Config{Address: ragCfg.Address})
		if err == nil {
			embedder, embedderErr := rag.NewConfiguredEmbedder(ragCfg)
			if embedderErr != nil {
				if logger != nil {
					logger.Warn("语义嵌入器初始化失败，流式服务降级为本地关键词检索", zap.Error(embedderErr))
				}
				return rag.NewLocalRetriever(rag.DefaultDocuments())
			}
			if logger != nil {
				logger.Info("AI 智能体流式服务已连接 Milvus", zap.String("address", ragCfg.Address))
			}
			return rag.NewMilvusRetriever(milvusClient, embedder)
		}
		if logger != nil {
			logger.Warn("无法连接 Milvus，流式服务降级为本地关键词检索", zap.String("address", ragCfg.Address), zap.Error(err))
		}
	}

	if logger != nil {
		logger.Info("AI 智能体流式服务使用本地关键词检索")
	}
	return rag.NewLocalRetriever(rag.DefaultDocuments())
}
