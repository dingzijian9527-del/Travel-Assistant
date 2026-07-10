package rpcaiagent

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/rag"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/traveldata"
)

func TestAIAgentServiceChatUsesModelClient(t *testing.T) {
	model := &fakeModelClient{reply: "上海住宿建议：优先选人民广场、静安寺或南京西路附近，兼顾地铁和餐饮。"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model, nil, config.RAGConfig{})
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
	model := &fakeModelClient{reply: "继续沿用成都三天行程，第二天建议增加人民公园茶坊。"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model, nil, config.RAGConfig{})
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
	if len(model.lastMessages) < 3 || !modelMessagesContain(model.lastMessages, "成都三天怎么玩") {
		t.Fatalf("expected previous conversation history, got %#v", model.lastMessages)
	}
}

func TestAIAgentServiceAddsTravelDataContextToModelMessages(t *testing.T) {
	model := &fakeModelClient{reply: "成都三天行程可以按天气、路线和预算安排。"}
	service := newAIAgentServiceWithTravelData(
		newAIAgentRepo(),
		model,
		nil,
		config.RAGConfig{},
		fakeTravelPlanner{},
	)

	_, _, svcErr := service.Chat(context.Background(), &aiagent.ChatRequest{
		UserId:  8,
		Message: "帮我规划成都三天行程，预算3000元",
	})
	if svcErr != nil {
		t.Fatalf("chat returned error: %v", svcErr)
	}
	for _, expected := range []string{"旅行实时参考", "成都", "多云", "宽窄巷子", "预算拆分", "路线参考"} {
		if !modelMessagesContain(model.lastMessages, expected) {
			t.Fatalf("expected model messages to contain %s, got %#v", expected, model.lastMessages)
		}
	}
}

func TestBuildModelMessagesAddsConversationStateRules(t *testing.T) {
	messages := buildModelMessages(nil, "想出去玩", "")
	if len(messages) == 0 {
		t.Fatal("expected model messages")
	}
	if !strings.Contains(messages[0].Content, "对话状态规则") {
		t.Fatalf("expected conversation state rules, got %#v", messages[0])
	}
	if !strings.Contains(messages[0].Content, "先追问缺失信息") {
		t.Fatalf("expected missing information rule, got %#v", messages[0])
	}
	if !strings.Contains(messages[0].Content, "热门推荐") {
		t.Fatalf("expected popular recommendation fallback rule, got %#v", messages[0])
	}
}

func TestBuildModelMessagesKeepsHistoryForPartialFollowUp(t *testing.T) {
	history := []chatMessageModel{
		{Role: aiagent.ChatRole_USER, Content: "我想去哈尔滨旅游"},
		{Role: aiagent.ChatRole_ASSISTANT, Content: "你计划什么时候去哈尔滨？"},
	}
	messages := buildModelMessages(history, "7月，三天", "")
	if !modelMessagesContain(messages, "我想去哈尔滨旅游") {
		t.Fatalf("expected destination history, got %#v", messages)
	}
	if !modelMessagesContain(messages, "7月，三天") {
		t.Fatalf("expected current partial follow-up, got %#v", messages)
	}
	if !strings.Contains(messages[0].Content, "用户后续只补充") {
		t.Fatalf("expected partial follow-up rule, got %#v", messages[0])
	}
}

func TestAIAgentServiceGuidesNonTravelMessage(t *testing.T) {
	model := &fakeModelClient{reply: "不应该调用"}
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model, nil, config.RAGConfig{})
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

func TestAIAgentServiceChatStreamWithRetriever(t *testing.T) {
	model := &fakeModelClient{reply: "成都美食推荐火锅、串串和担担面。"}
	documents := rag.DefaultDocuments()
	retriever := rag.NewLocalRetriever(documents)
	service := newAIAgentServiceWithModel(newAIAgentRepo(), model, retriever, config.RAGConfig{
		Enabled:  true,
		TopK:     3,
		MinScore: 0.1,
	})

	var buf bytes.Buffer
	svcErr := service.ChatStream(context.Background(), 2, "成都有什么好吃的", &buf)
	if svcErr != nil {
		t.Fatalf("chat stream returned error: %v", svcErr)
	}
	if buf.Len() == 0 {
		t.Fatal("expected stream output")
	}
}

func TestAIAgentServiceChatStreamEmptyMessage(t *testing.T) {
	service := newAIAgentService(newAIAgentRepo())
	svcErr := service.ChatStream(context.Background(), 1, "", io.Discard)
	if svcErr == nil {
		t.Fatal("expected validation error for empty message")
	}
}

type fakeModelClient struct {
	reply        string
	called       bool
	lastMessages []modelMessage
}

type fakeTravelPlanner struct{}

func (fakeTravelPlanner) BuildContext(ctx context.Context, req traveldata.Request) (traveldata.Result, error) {
	return traveldata.Result{
		Destination: req.Destination,
		Days:        req.Days,
		People:      req.People,
		Weather: []traveldata.WeatherDay{
			{Date: "2026-07-06", Text: "多云", TempMin: 24, TempMax: 32, Wind: "微风"},
		},
		Places: []traveldata.Place{
			{Name: "宽窄巷子", Address: "青羊区", Category: "景点"},
			{Name: "春熙路", Address: "锦江区", Category: "美食"},
		},
		Routes: []traveldata.Route{
			{From: "宽窄巷子", To: "春熙路", DistanceMeters: 3200, DurationMinutes: 18},
		},
		BudgetItems: []traveldata.BudgetItem{
			{Label: "住宿", Amount: 1200},
			{Label: "餐饮", Amount: 750},
		},
	}, nil
}

func (f *fakeModelClient) Chat(ctx context.Context, messages []modelMessage) (string, error) {
	f.called = true
	f.lastMessages = append([]modelMessage(nil), messages...)
	return f.reply, nil
}

func (f *fakeModelClient) StreamChat(ctx context.Context, messages []modelMessage, output io.Writer) error {
	f.called = true
	f.lastMessages = append([]modelMessage(nil), messages...)
	_, err := output.Write([]byte(f.reply))
	return err
}

func modelMessagesContain(messages []modelMessage, content string) bool {
	for _, message := range messages {
		if strings.Contains(message.Content, content) {
			return true
		}
	}
	return false
}
