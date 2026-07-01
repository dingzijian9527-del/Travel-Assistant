package rag

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// Retriever searches travel knowledge for the current user question.
type Retriever interface {
	Search(ctx context.Context, query string, cfg config.RAGConfig) ([]Result, error)
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
