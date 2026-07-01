package rpcaiagent

import "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"

type serviceError struct {
	code    base.ErrorCode
	message string
}

func errParam(message string) *serviceError {
	return &serviceError{code: base.ErrorCode_PARAM_ERROR, message: message}
}
