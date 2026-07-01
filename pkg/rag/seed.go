package rag

import "strings"

// SeedRecord 是准备写入向量数据库的一行资料。
type SeedRecord struct {
	ID        string
	Title     string
	City      string
	Tags      string
	Content   string
	Embedding []float32
}

// BuildSeedRecords 把资料文档转换成向量数据库记录。
func BuildSeedRecords(documents []Document, embedder Embedder) []SeedRecord {
	records := make([]SeedRecord, 0, len(documents))
	for _, doc := range documents {
		records = append(records, SeedRecord{
			ID:        doc.ID,
			Title:     doc.Title,
			City:      doc.City,
			Tags:      strings.Join(doc.Tags, ","),
			Content:   doc.Content,
			Embedding: embedder.Embed(DocumentEmbeddingText(doc)),
		})
	}
	return records
}

// DocumentEmbeddingText 返回用于生成向量的可检索文本。
func DocumentEmbeddingText(doc Document) string {
	parts := []string{doc.Title, doc.City, strings.Join(doc.Tags, " "), doc.Content}
	return strings.Join(parts, "\n")
}
