namespace go base

// ErrorCode 定义全局共享的业务错误分类。
enum ErrorCode {
    // SUCCESS 表示请求处理成功。
    SUCCESS = 0,
    // INTERNAL_ERROR 表示服务内部发生非预期错误。
    INTERNAL_ERROR = 10001,
    // PARAM_ERROR 表示请求参数不合法。
    PARAM_ERROR = 10002,
    // AUTH_ERROR 表示认证或授权失败。
    AUTH_ERROR = 10003,
    // BIZ_ERROR 表示业务规则校验失败。
    BIZ_ERROR = 10004,
    // NOT_FOUND 表示请求的资源不存在。
    NOT_FOUND = 10005,
}

// BaseResp 是远程过程调用接口的统一响应包装。
struct BaseResp {
    1: required i32 code,
    2: required string msg,
    3: optional string traceId,
}
