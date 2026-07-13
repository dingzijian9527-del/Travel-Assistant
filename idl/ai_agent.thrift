namespace go ai_agent

include "base.thrift"

// ChatRole 定义对话消息角色。
enum ChatRole {
    // USER 表示用户发送的消息。
    USER = 1,
    // ASSISTANT 表示旅行智能体生成的消息。
    ASSISTANT = 2,
}

// ChatMessage 表示旅行智能体对话中的单条消息。
struct ChatMessage {
    1: required ChatRole role,
    2: required string content,
}

// ChatRequest 是旅行智能体单轮对话请求。
struct ChatRequest {
    1: required i64 userId,
    2: required string message,
    3: optional list<ChatMessage> history,
    4: optional string traceId,
}

// ChatResponse 是旅行智能体单轮对话响应。
struct ChatResponse {
    1: required base.BaseResp baseResp,
    2: optional ChatMessage reply,
    3: optional list<string> suggestions,
}

// PromptSuggestionsRequest 是旅行提示词推荐请求。
struct PromptSuggestionsRequest {
    1: optional i64 userId,
    2: optional string scene,
}

// PromptSuggestionsResponse 是旅行提示词推荐响应。
struct PromptSuggestionsResponse {
    1: required base.BaseResp baseResp,
    2: optional list<string> suggestions,
}

struct ChatStreamChunk {
    1: required base.BaseResp baseResp,
    2: optional string content,
    3: required bool done,
}

// AIAgentService 提供旅行智能体能力。
service AIAgentService {
    ChatResponse Chat(1: ChatRequest req),
    ChatStreamChunk ChatStream(1: ChatRequest req) (streaming.mode="server"),
    PromptSuggestionsResponse GetPromptSuggestions(1: PromptSuggestionsRequest req),
}
