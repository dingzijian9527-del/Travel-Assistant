namespace go user

include "base.thrift"

// RegisterRequest 鏄敤鎴锋敞鍐岃姹傘€?
struct RegisterRequest {
    1: required string phone,
    2: required string password,
    3: optional string nickname,
    4: optional string avatarUrl,
    5: optional string homeCity,
    6: optional string currentCity,
}

// UserInfo 鎻忚堪鐢ㄦ埛鍏紑璧勬枡瀛楁銆?
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

// RegisterResponse 鏄敤鎴锋敞鍐屽搷搴斻€?
struct RegisterResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

// LoginRequest 鏄敤鎴风櫥褰曡姹傘€?
struct LoginRequest {
    1: required string phone,
    2: required string password,
}

// LoginResponse 鏄敤鎴风櫥褰曞搷搴斻€?
struct LoginResponse {
    1: required base.BaseResp baseResp,
    2: optional string token,
    3: optional UserInfo user,
}

// GetProfileRequest 鏄煡璇㈢敤鎴疯祫鏂欒姹傘€?
struct GetProfileRequest {
    1: required string id,
}

// GetProfileResponse 鏄煡璇㈢敤鎴疯祫鏂欏搷搴斻€?
struct GetProfileResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

// UpdateProfileRequest 鏄洿鏂扮敤鎴疯祫鏂欒姹傘€?
struct UpdateProfileRequest {
    1: required string id,
    2: optional string nickname,
    3: optional string avatarUrl,
    4: optional string homeCity,
    5: optional string currentCity,
}

// UpdateProfileResponse 鏄洿鏂扮敤鎴疯祫鏂欏搷搴斻€?
struct UpdateProfileResponse {
    1: required base.BaseResp baseResp,
    2: optional UserInfo user,
}

struct UserStats {
    1: required i64 tripCount,
    2: required i64 favoriteCount,
    3: required i64 unreadCount,
    4: required i64 couponCount,
}

struct UserSettings {
    1: required bool tripReminderEnabled,
    2: required bool priceReminderEnabled,
    3: required bool personalizedRecommendEnabled,
}

struct UserPreferences {
    1: required list<string> items,
}

struct UserDashboard {
    1: optional UserInfo user,
    2: optional UserStats stats,
    3: optional UserSettings settings,
    4: optional UserPreferences preferences,
}

struct GetDashboardRequest {
    1: required string id,
}

struct GetDashboardResponse {
    1: required base.BaseResp baseResp,
    2: optional UserDashboard dashboard,
}

struct GetPreferencesRequest {
    1: required string id,
}

struct GetPreferencesResponse {
    1: required base.BaseResp baseResp,
    2: optional UserPreferences preferences,
}

struct UpdatePreferencesRequest {
    1: required string id,
    2: required list<string> items,
}

struct UpdatePreferencesResponse {
    1: required base.BaseResp baseResp,
    2: optional UserPreferences preferences,
}

struct GetSettingsRequest {
    1: required string id,
}

struct GetSettingsResponse {
    1: required base.BaseResp baseResp,
    2: optional UserSettings settings,
}

struct UpdateSettingsRequest {
    1: required string id,
    2: required bool tripReminderEnabled,
    3: required bool priceReminderEnabled,
    4: required bool personalizedRecommendEnabled,
}

struct UpdateSettingsResponse {
    1: required base.BaseResp baseResp,
    2: optional UserSettings settings,
}

// UserService 鎻愪緵鐢ㄦ埛璐﹀彿涓庤祫鏂欒兘鍔涖€?
service UserService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    GetProfileResponse GetProfile(1: GetProfileRequest req),
    UpdateProfileResponse UpdateProfile(1: UpdateProfileRequest req),
    GetDashboardResponse GetDashboard(1: GetDashboardRequest req),
    GetPreferencesResponse GetPreferences(1: GetPreferencesRequest req),
    UpdatePreferencesResponse UpdatePreferences(1: UpdatePreferencesRequest req),
    GetSettingsResponse GetSettings(1: GetSettingsRequest req),
    UpdateSettingsResponse UpdateSettings(1: UpdateSettingsRequest req),
}
