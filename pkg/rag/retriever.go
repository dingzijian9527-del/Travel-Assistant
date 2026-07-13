package rag

import (
	"context"
	"errors"
	"sort"
	"strings"
	"unicode"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// Retriever searches travel knowledge for the current user question.
type Retriever interface {
	Search(ctx context.Context, query string, cfg config.RAGConfig) ([]Result, error)
}

// MilvusRetriever 从 Milvus 向量数据库检索旅行知识。
type MilvusRetriever struct {
	milvusClient client.Client
	embedder     Embedder
}

// NewMilvusRetriever 创建基于 Milvus 的向量检索器。
func NewMilvusRetriever(milvusClient client.Client, embedder Embedder) *MilvusRetriever {
	return &MilvusRetriever{milvusClient: milvusClient, embedder: embedder}
}

// Search 将用户问题转成向量后在 Milvus 中做近似搜索，返回最相关的知识条目。
func (r *MilvusRetriever) Search(ctx context.Context, query string, cfg config.RAGConfig) ([]Result, error) {
	if !cfg.Enabled || strings.TrimSpace(query) == "" || cfg.CollectionName == "" {
		return nil, nil
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	queryEmbedding, err := embedQuery(ctx, r.embedder, query)
	if err != nil {
		return nil, err
	}
	queryVector := entity.FloatVector(queryEmbedding)
	searchParams, _ := entity.NewIndexFlatSearchParam()

	searchResult, err := r.milvusClient.Search(
		ctx,
		cfg.CollectionName,
		nil,
		"",
		[]string{"id", "title", "city", "tags", "content"},
		[]entity.Vector{queryVector},
		"embedding",
		entity.L2,
		cfg.TopK,
		searchParams,
	)
	if err != nil {
		return nil, err
	}
	if len(searchResult) == 0 {
		return nil, nil
	}

	results := make([]Result, 0, len(searchResult))
	for _, sr := range searchResult {
		fields := sr.Fields
		if len(fields) < 5 {
			continue
		}
		// 取该结果的第一条匹配记录（SearchResult 可能包含多条）
		results = append(results, extractResults(fields, sr.Scores, cfg.MinScore)...)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Document.ID < results[j].Document.ID
		}
		return results[i].Score > results[j].Score
	})
	if cfg.TopK > 0 && len(results) > cfg.TopK {
		results = results[:cfg.TopK]
	}
	return results, nil
}

func embedQuery(ctx context.Context, embedder Embedder, query string) ([]float32, error) {
	if contextual, ok := embedder.(ContextEmbedder); ok {
		return contextual.EmbedContext(ctx, query)
	}
	vector := embedder.Embed(query)
	if len(vector) == 0 {
		return nil, errors.New("嵌入器没有返回向量")
	}
	return vector, nil
}

func extractResults(fields client.ResultSet, scores []float32, minScore float64) []Result {
	if len(fields) < 5 || len(scores) == 0 {
		return nil
	}
	count := fields[0].Len()
	if count > len(scores) {
		count = len(scores)
	}
	results := make([]Result, 0, count)
	for i := 0; i < count; i++ {
		result := Result{Score: float64(scores[i])}
		if score := float64(scores[i]); score < minScore {
			continue
		}
		if idCol, ok := fields[0].(*entity.ColumnVarChar); ok {
			if val, err := idCol.ValueByIdx(i); err == nil {
				result.Document.ID = val
			}
		}
		if titleCol, ok := fields[1].(*entity.ColumnVarChar); ok {
			if val, err := titleCol.ValueByIdx(i); err == nil {
				result.Document.Title = val
			}
		}
		if cityCol, ok := fields[2].(*entity.ColumnVarChar); ok {
			if val, err := cityCol.ValueByIdx(i); err == nil {
				result.Document.City = val
			}
		}
		if tagsCol, ok := fields[3].(*entity.ColumnVarChar); ok {
			if val, err := tagsCol.ValueByIdx(i); err == nil {
				result.Document.Tags = splitComma(val)
			}
		}
		if contentCol, ok := fields[4].(*entity.ColumnVarChar); ok {
			if val, err := contentCol.ValueByIdx(i); err == nil {
				result.Document.Content = val
			}
		}
		results = append(results, result)
	}
	return results
}

func splitComma(val string) []string {
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// LocalRetriever is a deterministic keyword retriever for built-in snippets.
type LocalRetriever struct {
	documents []Document
}

// NewLocalRetriever creates a local retriever from documents.
func NewLocalRetriever(documents []Document) *LocalRetriever {
	copied := make([]Document, len(documents))
	copy(copied, documents)
	return &LocalRetriever{documents: copied}
}

// Search returns documents ranked by simple keyword overlap.
func (r *LocalRetriever) Search(ctx context.Context, query string, cfg config.RAGConfig) ([]Result, error) {
	if !cfg.Enabled || strings.TrimSpace(query) == "" {
		return nil, nil
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	queryTokens := tokenize(query)
	if len(queryTokens) == 0 {
		return nil, nil
	}

	results := make([]Result, 0, len(r.documents))
	for _, doc := range r.documents {
		score := scoreDocument(queryTokens, doc)
		if score >= cfg.MinScore {
			results = append(results, Result{Document: doc, Score: score})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Document.ID < results[j].Document.ID
		}
		return results[i].Score > results[j].Score
	})
	if cfg.TopK > 0 && len(results) > cfg.TopK {
		results = results[:cfg.TopK]
	}
	return results, nil
}

func scoreDocument(queryTokens map[string]struct{}, doc Document) float64 {
	target := doc.Title + " " + doc.City + " " + strings.Join(doc.Tags, " ") + " " + doc.Content
	docTokens := tokenize(target)
	if len(docTokens) == 0 {
		return 0
	}
	matches := 0
	for token := range queryTokens {
		if _, ok := docTokens[token]; ok {
			matches++
		}
	}
	return float64(matches) / float64(len(queryTokens))
}

func tokenize(text string) map[string]struct{} {
	tokens := make(map[string]struct{})
	var current strings.Builder
	flush := func() {
		if current.Len() == 0 {
			return
		}
		tokens[strings.ToLower(current.String())] = struct{}{}
		current.Reset()
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
			continue
		}
		flush()
	}
	flush()

	for _, keyword := range travelKeywords {
		if strings.Contains(text, keyword) {
			tokens[keyword] = struct{}{}
		}
	}
	return tokens
}

var travelKeywords = []string{
	"成都", "杭州", "北京", "上海", "广州", "三亚",
	"美食", "住宿", "酒店", "行程", "路线", "交通", "亲子", "周末",
	"西湖", "火锅", "早茶", "海岛", "度假", "攻略", "预算", "景点",
}
