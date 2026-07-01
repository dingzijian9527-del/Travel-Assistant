package rag

import (
	"context"
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestLocalRetrieverFindsChengduKnowledge(t *testing.T) {
	retriever := NewLocalRetriever(DefaultDocuments())

	results, err := retriever.Search(context.Background(), "成都三天美食和住宿怎么安排", config.RAGConfig{
		Enabled:  true,
		TopK:     3,
		MinScore: 0.1,
	})
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected chengdu result")
	}
	if results[0].Document.City != "成都" {
		t.Fatalf("expected chengdu as top result, got %#v", results[0])
	}
}

func TestLocalRetrieverReturnsEmptyForLowScore(t *testing.T) {
	retriever := NewLocalRetriever(DefaultDocuments())

	results, err := retriever.Search(context.Background(), "量子力学作业怎么写", config.RAGConfig{
		Enabled:  true,
		TopK:     3,
		MinScore: 0.2,
	})
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %#v", results)
	}
}
