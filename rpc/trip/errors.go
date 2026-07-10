package rpctrip

import "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"

type serviceError struct {
	code    base.ErrorCode
	message string
}

func (e *serviceError) Error() string {
	return e.message
}

func errParam(message string) *serviceError {
	return &serviceError{code: base.ErrorCode_PARAM_ERROR, message: message}
}

func errBiz(message string) *serviceError {
	return &serviceError{code: base.ErrorCode_BIZ_ERROR, message: message}
}

func errInternal(message string) *serviceError {
	return &serviceError{code: base.ErrorCode_INTERNAL_ERROR, message: message}
}

func errNotFound(message string) *serviceError {
	return &serviceError{code: base.ErrorCode_NOT_FOUND, message: message}
}
