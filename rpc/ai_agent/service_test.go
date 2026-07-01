package rpcaiagent

import (
	"context"
	"strings"
	"testing"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
)

func TestAIAgentServiceChatUsesModelClient(t *testing.T) {
	model := &fakeModelClient{reply: "上海住宿建议：优先选人民广场、静安寺或南京西路附近，兼顾地铁和餐饮。"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model)
	reply, suggestions, svcErr := service.Chat(context.Background(), &aiagent.ChatRequest{
		UserId:  1,
		Message: "帮我搜索上海精品酒店",
	})
	if svcErr != nil {
		t.Fatalf("chat returned error: %v", svcErr)
	}
	if !model.called {
		t.Fatal("expected model client to be called")
	}
	if reply == nil || reply.Content != model.reply {
		t.Fatalf("unexpected reply: %+v", reply)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected suggestions")
	}
}

func TestAIAgentServicePassesHistoryToModelClient(t *testing.T) {
	model := &fakeModelClient{reply: "继续沿用成都三天行程，第二天建议增加人民公园茶馆。"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model)
	reply, _, svcErr := service.Chat(context.Background(), &aiagent.ChatRequest{
		UserId:  1,
		Message: "成都三天怎么玩",
	})
	if svcErr != nil {
		t.Fatalf("first chat returned error: %v", svcErr)
	}
	reply, _, svcErr = service.Chat(context.Background(), &aiagent.ChatRequest{
		UserId:  1,
		Message: "第二天不要太赶",
	})
	if svcErr != nil {
		t.Fatalf("chat returned error: %v", svcErr)
	}
	if reply == nil || !strings.Contains(reply.Content, "人民公园") {
		t.Fatalf("unexpected reply: %+v", reply)
	}
	if len(model.lastMessages) < 2 || !strings.Contains(model.lastMessages[0].Content, "成都三天怎么玩") {
		t.Fatalf("expected previous conversation history, got %#v", model.lastMessages)
	}
}

func TestAIAgentServiceGuidesNonTravelMessage(t *testing.T) {
	model := &fakeModelClient{reply: "不应该调用"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model)
	reply, suggestions, svcErr := service.Chat(context.Background(), &aiagent.ChatRequest{
		UserId:  1,
		Message: "帮我写一个排序算法",
	})
	if svcErr != nil {
		t.Fatalf("chat returned error: %v", svcErr)
	}
	if reply == nil || !strings.Contains(reply.Content, "旅游专用") {
		t.Fatalf("unexpected non-travel reply: %+v", reply)
	}
	if len(suggestions) == 0 || !strings.Contains(suggestions[0], "行程") {
		t.Fatalf("unexpected suggestions: %+v", suggestions)
	}
	if model.called {
		t.Fatal("non-travel message should not call model client")
	}
}

func TestAIAgentServiceRejectsEmptyMessage(t *testing.T) {
	service := newAIAgentService(newAIAgentRepo())
	_, _, svcErr := service.Chat(context.Background(), &aiagent.ChatRequest{UserId: 1})
	if svcErr == nil {
		t.Fatal("expected validation error")
	}
}

type fakeModelClient struct {
	reply        string
	called       bool
	lastMessages []modelMessage
}

func (f *fakeModelClient) Chat(ctx context.Context, messages []modelMessage) (string, error) {
	f.called = true
	f.lastMessages = append([]modelMessage(nil), messages...)
	return f.reply, nil
}
