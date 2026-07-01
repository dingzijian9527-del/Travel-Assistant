package rpcaiagent

import (
	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"
)

func toChatMessageDTO(input *chatMessageModel) *aiagent.ChatMessage {
	if input == nil {
		return nil
	}
	return &aiagent.ChatMessage{Role: input.Role, Content: input.Content}
}

func successBaseResp() *base.BaseResp {
	return &base.BaseResp{Code: int32(base.ErrorCode_SUCCESS), Msg: "success"}
}

func errorBaseResp(err *serviceError) *base.BaseResp {
	if err == nil {
		return successBaseResp()
	}
	return &base.BaseResp{Code: int32(err.code), Msg: err.message}
}

func chatErrorResponse(err *serviceError) *aiagent.ChatResponse {
	return &aiagent.ChatResponse{BaseResp: errorBaseResp(err)}
}
