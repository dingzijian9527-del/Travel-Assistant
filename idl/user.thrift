namespace go user

include "base.thrift"

// RegisterRequest 是用户注册请求。
struct RegisterRequest {
    1: required string phone,
    2: required string password,
    3: optional string nickname,
    4: optional string avatarUrl,
    5: optional string homeCity,
    6: optional string currentCity,
}

// UserInfo 描述用户公开资料字段。
struct UserInfo {
    1: required string id,
    2: required string phone,
    3: required string nickname,
    4: optional string avatarUrl,
    5: optional string homeCity,
    6: optional string currentCity,
    7: required string memberLevel,
    8: required i32 accountStatus,
    9: optional string createdAt,
    10: optional string updatedAt,
    11: optional string deletedAt,
}

// RegisterResponse 是用户注册响应。
struct RegisterResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

// LoginRequest 是用户登录请求。
struct LoginRequest {
    1: required string phone,
    2: required string password,
}

// LoginResponse 是用户登录响应。
struct LoginResponse {
    1: required base.BaseResp baseResp,
    2: optional string token,
    3: optional UserInfo user,
}

// GetProfileRequest 是查询用户资料请求。
struct GetProfileRequest {
    1: required string id,
}

// GetProfileResponse 是查询用户资料响应。
struct GetProfileResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

// UpdateProfileRequest 是更新用户资料请求。
struct UpdateProfileRequest {
    1: required string id,
    2: optional string nickname,
    3: optional string avatarUrl,
    4: optional string homeCity,
    5: optional string currentCity,
}

// UpdateProfileResponse 是更新用户资料响应。
struct UpdateProfileResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

// UserService 提供用户账号与资料能力。
service UserService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    GetProfileResponse GetProfile(1: GetProfileRequest req),
    UpdateProfileResponse UpdateProfile(1: UpdateProfileRequest req),
}
