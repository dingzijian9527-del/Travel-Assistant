package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func main() {
	configPath := flag.String("config", "conf/config.yaml", "配置文件路径")
	reset := flag.Bool("reset", false, "写入前是否删除并重建集合")
	timeout := flag.Duration("timeout", 30*time.Second, "连接和写入超时时间")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}
	if cfg.RAG.Address == "" {
		log.Fatal("Milvus 地址不能为空，请检查 rag.address")
	}
	if cfg.RAG.CollectionName == "" {
		log.Fatal("集合名称不能为空，请检查 rag.collection_name")
	}
	if cfg.RAG.EmbeddingDim <= 0 {
		log.Fatalf("向量维度必须大于 0，当前值：%d", cfg.RAG.EmbeddingDim)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	milvusClient, err := client.NewClient(ctx, client.Config{Address: cfg.RAG.Address})
	if err != nil {
		log.Fatalf("连接 Milvus 失败：%v", err)
	}
	defer milvusClient.Close()

	documents := rag.DefaultDocuments()
	records := rag.BuildSeedRecords(documents, rag.NewHashEmbedder(cfg.RAG.EmbeddingDim))
	writer := rag.NewMilvusWriter(milvusClient)

	if err := writer.EnsureCollection(ctx, cfg.RAG.CollectionName, cfg.RAG.EmbeddingDim, *reset); err != nil {
		log.Fatalf("准备集合失败：%v", err)
	}
	if err := writer.UpsertRecords(ctx, cfg.RAG.CollectionName, records, cfg.RAG.EmbeddingDim); err != nil {
		log.Fatalf("写入资料失败：%v", err)
	}

	fmt.Printf("已写入 %d 条资料到集合 %s\n", len(records), cfg.RAG.CollectionName)
}
