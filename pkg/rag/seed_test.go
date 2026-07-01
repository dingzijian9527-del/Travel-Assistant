package rag

import (
	"strings"
	"testing"
)

func TestBuildSeedRecordsUsesDocumentFieldsAndEmbedding(t *testing.T) {
	docs := []Document{
		{
			ID:      "chengdu-3-days",
			Title:   "成都三日游行程与住宿",
			City:    "成都",
			Tags:    []string{"行程", "美食"},
			Content: "适合住春熙路、太古里或宽窄巷子附近。",
		},
	}

	records := BuildSeedRecords(docs, NewHashEmbedder(6))

	if len(records) != 1 {
		t.Fatalf("expected one seed record, got %d", len(records))
	}
	record := records[0]
	if record.ID != "chengdu-3-days" {
		t.Fatalf("unexpected id: %s", record.ID)
	}
	if record.Title != "成都三日游行程与住宿" {
		t.Fatalf("unexpected title: %s", record.Title)
	}
	if record.City != "成都" {
		t.Fatalf("unexpected city: %s", record.City)
	}
	if record.Tags != "行程,美食" {
		t.Fatalf("unexpected tags: %s", record.Tags)
	}
	if record.Content != "适合住春熙路、太古里或宽窄巷子附近。" {
		t.Fatalf("unexpected content: %s", record.Content)
	}
	if len(record.Embedding) != 6 {
		t.Fatalf("expected embedding dimension 6, got %d", len(record.Embedding))
	}
}

func TestDocumentEmbeddingTextIncludesSearchableMetadata(t *testing.T) {
	doc := Document{
		Title:   "广州美食与老城路线",
		City:    "广州",
		Tags:    []string{"美食", "早茶"},
		Content: "老城可逛陈家祠、沙面、永庆坊。",
	}

	text := DocumentEmbeddingText(doc)

	for _, want := range []string{"广州美食与老城路线", "广州", "美食", "早茶", "陈家祠"} {
		if !strings.Contains(text, want) {
			t.Fatalf("embedding text should contain %q, got %s", want, text)
		}
	}
}
