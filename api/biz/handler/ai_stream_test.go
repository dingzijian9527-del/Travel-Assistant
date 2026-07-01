package handler

import (
	"context"
	"io"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestBuildAIContextMessagesKeepsConversationHistory(t *testing.T) {
	history := []conversationMessage{
		{Role: "user", Content: "我要去上海玩7天"},
		{Role: "assistant", Content: "可以，我先按上海7天帮你排。"},
		{Role: "user", Content: "天气如何"},
	}

	messages := buildAIContextMessages(history, "那住宿怎么选")

	if len(messages) != 4 {
		t.Fatalf("expected history plus current message, got %#v", messages)
	}
	if messages[0].Role != "user" || messages[1].Role != "assistant" || messages[2].Content != "天气如何" {
		t.Fatalf("history messages not preserved: %#v", messages)
	}
	if messages[3].Role != "user" || messages[3].Content != "那住宿怎么选" {
		t.Fatalf("current user message not appended: %#v", messages)
	}
}

func TestBuildAIContextMessagesAddsRAGContext(t *testing.T) {
	messages := buildAIContextMessages(nil, "成都三天怎么玩", "可参考资料：成都适合住春熙路附近")

	if len(messages) != 1 {
		t.Fatalf("expected one message, got %#v", messages)
	}
	if messages[0].Role != "user" {
		t.Fatalf("expected user message, got %#v", messages[0])
	}
	if !strings.Contains(messages[0].Content, "可参考资料") {
		t.Fatalf("rag context missing: %s", messages[0].Content)
	}
	if !strings.Contains(messages[0].Content, "用户问题：成都三天怎么玩") {
		t.Fatalf("user question missing: %s", messages[0].Content)
	}
}

func TestBuildAIContextMessagesSkipsEmptyRAGContext(t *testing.T) {
	messages := buildAIContextMessages(nil, "杭州两天怎么玩", "  ")

	if len(messages) != 1 {
		t.Fatalf("expected one message, got %#v", messages)
	}
	if messages[0].Content != "杭州两天怎么玩" {
		t.Fatalf("empty rag context should keep original message, got %s", messages[0].Content)
	}
}

func TestUnavailableModelDoesNotUseLocalTravelTemplate(t *testing.T) {
	travelSessions = newTravelSessionStore(12)
	runtime := &bootstrap.Runtime{
		Config: &config.Config{AI: config.AIConfig{MaxPromptChars: 2000}},
		Logger: zap.NewNop(),
	}
	reader, writer := io.Pipe()
	result := make(chan streamReadResult, 1)

	go func() {
		var output strings.Builder
		_, err := io.Copy(&output, reader)
		result <- streamReadResult{text: output.String(), err: err}
	}()
	writeSmartTravelReply(context.Background(), writer, runtime, 1001, "我要去上海玩7天帮我规划路线")

	got := <-result
	if got.err != nil {
		t.Fatalf("read stream failed: %v", got.err)
	}
	if !strings.Contains(got.text, "智能体模型暂不可用") {
		t.Fatalf("fallback should explain model unavailable, got %s", got.text)
	}
	for _, forbidden := range []string{"第1天", "上海7天", "路线建议", "本地特色餐"} {
		if strings.Contains(got.text, forbidden) {
			t.Fatalf("fallback should not generate local travel content %q, got %s", forbidden, got.text)
		}
	}
	if history := travelSessions.History(1001); len(history) != 0 {
		t.Fatalf("unavailable model reply should not be stored as model history: %#v", history)
	}
}

type streamReadResult struct {
	text string
	err  error
}
