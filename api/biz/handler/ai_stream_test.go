package handler

import "testing"

func TestModelUnavailableReplyUsesReadableChinese(t *testing.T) {
	expected := "智能体模型暂不可用，请稍后重试或检查火山方舟接入点配置。"
	if modelUnavailableReply != expected {
		t.Fatalf("降级提示不符合预期：%q", modelUnavailableReply)
	}
}

func TestWriteChatStreamReplyFallbackMessageUsesReadableChinese(t *testing.T) {
	expected := "智能体服务暂不可用，请稍后重试。"
	if fallback := unavailableServiceReply(); fallback != expected {
		t.Fatalf("服务降级提示不符合预期：%q", fallback)
	}
}
