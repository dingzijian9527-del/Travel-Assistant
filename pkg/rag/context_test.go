package rag

import (
	"strings"
	"testing"
)

func TestFormatContextIncludesReferenceMetadata(t *testing.T) {
	text := FormatContext([]Result{
		{
			Document: Document{
				Title:   "成都三日游",
				City:    "成都",
				Tags:    []string{"美食", "住宿"},
				Content: "春熙路和太古里适合首次到成都的游客。",
			},
			Score: 0.6,
		},
	})

	for _, want := range []string{"可参考资料", "成都三日游", "成都", "美食", "春熙路"} {
		if !strings.Contains(text, want) {
			t.Fatalf("formatted context should contain %q, got %s", want, text)
		}
	}
}

func TestFormatContextReturnsEmptyForNoResults(t *testing.T) {
	if got := FormatContext(nil); got != "" {
		t.Fatalf("empty results should format to empty string, got %s", got)
	}
}
