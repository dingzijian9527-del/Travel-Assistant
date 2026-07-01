package rag

import "testing"

func TestHashEmbedderReturnsStableVectorWithConfiguredDimension(t *testing.T) {
	embedder := NewHashEmbedder(8)

	first := embedder.Embed("成都三日游适合住春熙路")
	second := embedder.Embed("成都三日游适合住春熙路")

	if len(first) != 8 {
		t.Fatalf("expected 8 dimensions, got %d", len(first))
	}
	if len(second) != 8 {
		t.Fatalf("expected 8 dimensions, got %d", len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("expected stable vector, first=%v second=%v", first, second)
		}
	}
}

func TestHashEmbedderSeparatesDifferentText(t *testing.T) {
	embedder := NewHashEmbedder(16)

	first := embedder.Embed("成都火锅和春熙路住宿")
	second := embedder.Embed("三亚海岛度假和亲子酒店")

	if equalFloat32Slices(first, second) {
		t.Fatalf("expected different texts to produce different vectors")
	}
}

func TestHashEmbedderFallsBackToOneDimension(t *testing.T) {
	embedder := NewHashEmbedder(0)

	vector := embedder.Embed("杭州西湖")

	if len(vector) != 1 {
		t.Fatalf("expected fallback dimension 1, got %d", len(vector))
	}
}

func equalFloat32Slices(first, second []float32) bool {
	if len(first) != len(second) {
		return false
	}
	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}
	return true
}
