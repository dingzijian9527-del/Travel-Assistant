namespace go trip

include "base.thrift"

// TripSummary 描述行程卡片和概要栏需要展示的摘要信息。
struct TripSummary {
    1: optional string date,
    2: optional string days,
    3: optional string people,
    4: optional string budget,
}

// TripTip 描述单日行程中的提醒项。
struct TripTip {
    1: optional string icon,
    2: optional string title,
    3: optional string text,
}

// TripDay 描述单日计划。
struct TripDay {
    1: required i32 day,
    2: optional string title,
    3: optional string route,
    4: optional string food,
    5: optional string hotel,
    6: optional list<TripTip> tips,
    7: optional string weather,
}

// TripBudget 描述预算构成项。
struct TripBudget {
    1: optional string label,
    2: optional string amount,
}

// TripInfo 描述一个完整行程。
struct TripInfo {
    1: required string id,
    2: required i64 userId,
    3: required string title,
    4: optional string subtitle,
    5: optional string destination,
    6: optional string dateRange,
    7: optional i32 dayCount,
    8: optional string people,
    9: optional string budgetLevel,
    10: optional string sourceQuestion,
    11: optional string sourceReply,
    12: optional TripSummary summary,
    13: optional list<TripDay> days,
    14: optional list<TripBudget> budget,
    15: optional list<string> alerts,
    16: required bool saved,
    17: optional string createdAt,
    18: optional string updatedAt,
}

// CreateTripRequest 是创建行程请求。
struct CreateTripRequest {
    1: required i64 userId,
    2: optional string title,
    3: optional string subtitle,
    4: optional string destination,
    5: optional string dateRange,
    6: optional i32 dayCount,
    7: optional string people,
    8: optional string budgetLevel,
    9: optional string sourceQuestion,
    10: optional string sourceReply,
    11: optional TripSummary summary,
    12: optional list<TripDay> days,
    13: optional list<TripBudget> budget,
    14: optional list<string> alerts,
}

struct CreateTripResponse {
    1: required base.BaseResp baseResp,
    2: optional TripInfo trip,
}

struct ListTripsRequest {
    1: required i64 userId,
    2: optional i32 limit,
}

struct ListTripsResponse {
    1: required base.BaseResp baseResp,
    2: optional list<TripInfo> trips,
}

struct GetLatestTripRequest {
    1: required i64 userId,
}

struct GetLatestTripResponse {
    1: required base.BaseResp baseResp,
    2: optional TripInfo trip,
}

struct GetTripDetailRequest {
    1: required i64 userId,
    2: required string tripId,
}

struct GetTripDetailResponse {
    1: required base.BaseResp baseResp,
    2: optional TripInfo trip,
}

struct DeleteTripRequest {
    1: required i64 userId,
    2: required string tripId,
}

struct DeleteTripResponse {
    1: required base.BaseResp baseResp,
}

struct UpdateTripRequest {
    1: required i64 userId,
    2: required string tripId,
    3: optional string title,
    4: optional string subtitle,
    5: optional string destination,
    6: optional string dateRange,
    7: optional i32 dayCount,
    8: optional string people,
    9: optional string budgetLevel,
    10: optional list<TripDay> days,
}

struct UpdateTripResponse {
    1: required base.BaseResp baseResp,
    2: optional TripInfo trip,
}

// TripService 提供用户行程管理能力。
service TripService {
    CreateTripResponse CreateTrip(1: CreateTripRequest req),
    ListTripsResponse ListTrips(1: ListTripsRequest req),
    GetLatestTripResponse GetLatestTrip(1: GetLatestTripRequest req),
    GetTripDetailResponse GetTripDetail(1: GetTripDetailRequest req),
    DeleteTripResponse DeleteTrip(1: DeleteTripRequest req),
    UpdateTripResponse UpdateTrip(1: UpdateTripRequest req),
}
