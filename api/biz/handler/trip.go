package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dingzijian9527-del/Travel-Assistant/api/biz/middleware"
	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/bootstrap"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

type createTripRequest struct {
	Title          string              `json:"title"`
	Subtitle       string              `json:"subtitle"`
	Destination    string              `json:"destination"`
	DateRange      string              `json:"date_range"`
	DayCount       int32               `json:"day_count"`
	People         string              `json:"people"`
	BudgetLevel    string              `json:"budget_level"`
	SourceQuestion string              `json:"source_question"`
	SourceReply    string              `json:"source_reply"`
	Summary        tripSummaryPayload  `json:"summary"`
	Days           []tripDayPayload    `json:"days"`
	Budget         []tripBudgetPayload `json:"budget"`
	Alerts         []string            `json:"alerts"`
}

type tripSummaryPayload struct {
	Date   string `json:"date"`
	Days   string `json:"days"`
	People string `json:"people"`
	Budget string `json:"budget"`
}

type tripDayPayload struct {
	Day     int32            `json:"day"`
	Title   string           `json:"title"`
	Route   string           `json:"route"`
	Food    string           `json:"food"`
	Hotel   string           `json:"hotel"`
	Tips    []tripTipPayload `json:"tips"`
	Weather string           `json:"weather"`
}

type tripTipPayload struct {
	Icon  string `json:"icon"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type tripBudgetPayload struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
}

func CreateTrip(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		var req createTripRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		resp, err := clients.trip.CreateTrip(ctx, toCreateTripRPCRequest(userID, req))
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetTrip())
	}
}

func ListTrips(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		resp, err := clients.trip.ListTrips(ctx, &trip.ListTripsRequest{UserId: userID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetTrips())
	}
}

func GetLatestTrip(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		resp, err := clients.trip.GetLatestTrip(ctx, &trip.GetLatestTripRequest{UserId: userID})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetTrip())
	}
}

func GetTripDetail(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		resp, err := clients.trip.GetTripDetail(ctx, &trip.GetTripDetailRequest{UserId: userID, TripId: c.Param("tripID")})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetTrip())
	}
}

func DeleteTrip(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		resp, err := clients.trip.DeleteTrip(ctx, &trip.DeleteTripRequest{UserId: userID, TripId: c.Param("tripID")})
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), nil)
	}
}

type updateTripRequest struct {
	Title       *string          `json:"title,omitempty"`
	Subtitle    *string          `json:"subtitle,omitempty"`
	Destination *string          `json:"destination,omitempty"`
	DateRange   *string          `json:"date_range,omitempty"`
	DayCount    *int32           `json:"day_count,omitempty"`
	People      *string          `json:"people,omitempty"`
	BudgetLevel *string          `json:"budget_level,omitempty"`
	Days        []tripDayPayload `json:"days,omitempty"`
}

func UpdateTrip(runtime *bootstrap.Runtime) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := tripUserIDFromRequest(runtime, c)
		if !ok {
			return
		}
		var req updateTripRequest
		if err := c.BindAndValidate(&req); err != nil {
			writeJSON(c, consts.StatusBadRequest, 400, "请求参数错误", nil)
			return
		}
		clients, err := clientsFor(runtime)
		if err != nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务暂不可用", nil)
			return
		}
		rpcReq := &trip.UpdateTripRequest{
			UserId:      userID,
			TripId:      c.Param("tripID"),
			Title:       req.Title,
			Subtitle:    req.Subtitle,
			Destination: req.Destination,
			DateRange:   req.DateRange,
			DayCount:    req.DayCount,
			People:      req.People,
			BudgetLevel: req.BudgetLevel,
			Days:        toTripDaysRPC(req.Days),
		}
		resp, err := clients.trip.UpdateTrip(ctx, rpcReq)
		if err != nil || resp == nil {
			writeJSON(c, consts.StatusInternalServerError, 500, "行程服务调用失败", nil)
			return
		}
		writeRPCResponse(c, resp.GetBaseResp().GetCode(), resp.GetBaseResp().GetMsg(), resp.GetTrip())
	}
}

func tripUserIDFromRequest(runtime *bootstrap.Runtime, c *app.RequestContext) (int64, bool) {
	claims, ok := middleware.AuthClaims(c)
	if !ok {
		writeJSON(c, consts.StatusUnauthorized, 401, "登录状态无效，请重新登录", nil)
		return 0, false
	}
	userID, err := jwtx.UserIDInt64(claims.UserID)
	if err != nil {
		writeJSON(c, consts.StatusUnauthorized, 401, "登录状态无效，请重新登录", nil)
		return 0, false
	}
	return userID, true
}

func toCreateTripRPCRequest(userID int64, req createTripRequest) *trip.CreateTripRequest {
	return &trip.CreateTripRequest{
		UserId:         userID,
		Title:          stringPtrIfNotEmpty(req.Title),
		Subtitle:       stringPtrIfNotEmpty(req.Subtitle),
		Destination:    stringPtrIfNotEmpty(req.Destination),
		DateRange:      stringPtrIfNotEmpty(req.DateRange),
		DayCount:       int32PtrIfPositive(req.DayCount),
		People:         stringPtrIfNotEmpty(req.People),
		BudgetLevel:    stringPtrIfNotEmpty(req.BudgetLevel),
		SourceQuestion: stringPtrIfNotEmpty(req.SourceQuestion),
		SourceReply:    stringPtrIfNotEmpty(req.SourceReply),
		Summary:        toTripSummaryRPC(req.Summary),
		Days:           toTripDaysRPC(req.Days),
		Budget:         toTripBudgetRPC(req.Budget),
		Alerts:         append([]string(nil), req.Alerts...),
	}
}

func toTripSummaryRPC(input tripSummaryPayload) *trip.TripSummary {
	return &trip.TripSummary{
		Date:   stringPtrIfNotEmpty(input.Date),
		Days:   stringPtrIfNotEmpty(input.Days),
		People: stringPtrIfNotEmpty(input.People),
		Budget: stringPtrIfNotEmpty(input.Budget),
	}
}

func toTripDaysRPC(input []tripDayPayload) []*trip.TripDay {
	items := make([]*trip.TripDay, 0, len(input))
	for _, item := range input {
		items = append(items, &trip.TripDay{
			Day:     item.Day,
			Title:   stringPtrIfNotEmpty(item.Title),
			Route:   stringPtrIfNotEmpty(item.Route),
			Food:    stringPtrIfNotEmpty(item.Food),
			Hotel:   stringPtrIfNotEmpty(item.Hotel),
			Tips:    toTripTipsRPC(item.Tips),
			Weather: stringPtrIfNotEmpty(item.Weather),
		})
	}
	return items
}

func toTripTipsRPC(input []tripTipPayload) []*trip.TripTip {
	items := make([]*trip.TripTip, 0, len(input))
	for _, item := range input {
		items = append(items, &trip.TripTip{
			Icon:  stringPtrIfNotEmpty(item.Icon),
			Title: stringPtrIfNotEmpty(item.Title),
			Text:  stringPtrIfNotEmpty(item.Text),
		})
	}
	return items
}

func toTripBudgetRPC(input []tripBudgetPayload) []*trip.TripBudget {
	items := make([]*trip.TripBudget, 0, len(input))
	for _, item := range input {
		items = append(items, &trip.TripBudget{
			Label:  stringPtrIfNotEmpty(item.Label),
			Amount: stringPtrIfNotEmpty(item.Amount),
		})
	}
	return items
}

func stringPtrIfNotEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func int32PtrIfPositive(value int32) *int32 {
	if value <= 0 {
		return nil
	}
	return &value
}
