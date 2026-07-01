package rag

import (
	"fmt"
	"strings"
)

// FormatContext renders retrieval hits for model context.
func FormatContext(results []Result) string {
	if len(results) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("以下是旅行知识库中的可参考资料，请结合用户问题使用；如果资料与问题无关，可以忽略。\n")
	for i, result := range results {
		doc := result.Document
		builder.WriteString(fmt.Sprintf("\n资料%d：%s\n", i+1, doc.Title))
		builder.WriteString(fmt.Sprintf("城市：%s\n", doc.City))
		if len(doc.Tags) > 0 {
			builder.WriteString(fmt.Sprintf("标签：%s\n", strings.Join(doc.Tags, "、")))
		}
		builder.WriteString(fmt.Sprintf("内容：%s\n", doc.Content))
	}
	return strings.TrimSpace(builder.String())
}
