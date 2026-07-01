package rag

import (
	"testing"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func TestBuildMilvusColumnsKeepsRecordOrderAndFields(t *testing.T) {
	records := []SeedRecord{
		{
			ID:        "chengdu-3-days",
			Title:     "成都三日游行程与住宿",
			City:      "成都",
			Tags:      "行程,美食",
			Content:   "适合住春熙路。",
			Embedding: []float32{0.1, 0.2},
		},
	}

	columns := BuildMilvusColumns(records, 2)

	if len(columns) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(columns))
	}
	assertVarCharValue(t, columns[0], "id", "chengdu-3-days")
	assertVarCharValue(t, columns[1], "title", "成都三日游行程与住宿")
	assertVarCharValue(t, columns[2], "city", "成都")
	assertVarCharValue(t, columns[3], "tags", "行程,美食")
	assertVarCharValue(t, columns[4], "content", "适合住春熙路。")

	vectorColumn, ok := columns[5].(*entity.ColumnFloatVector)
	if !ok {
		t.Fatalf("expected embedding column to be float vector, got %T", columns[5])
	}
	if vectorColumn.Name() != "embedding" {
		t.Fatalf("unexpected vector column name: %s", vectorColumn.Name())
	}
	if vectorColumn.Len() != 1 {
		t.Fatalf("expected one vector row, got %d", vectorColumn.Len())
	}
}

func assertVarCharValue(t *testing.T, column entity.Column, name string, want string) {
	t.Helper()
	if column.Name() != name {
		t.Fatalf("expected column %s, got %s", name, column.Name())
	}
	varCharColumn, ok := column.(*entity.ColumnVarChar)
	if !ok {
		t.Fatalf("expected varchar column %s, got %T", name, column)
	}
	got, err := varCharColumn.ValueByIdx(0)
	if err != nil {
		t.Fatalf("read column %s: %v", name, err)
	}
	if got != want {
		t.Fatalf("expected %s=%s, got %s", name, want, got)
	}
}
