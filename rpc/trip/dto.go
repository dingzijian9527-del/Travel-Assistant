package rpctrip

import (
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"
	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
)

func successBaseResp() *base.BaseResp {
	return &base.BaseResp{Code: int32(base.ErrorCode_SUCCESS), Msg: "success"}
}

func errorBaseResp(err *serviceError) *base.BaseResp {
	return &base.BaseResp{Code: int32(err.code), Msg: err.message}
}

func createTripErrorResponse(err *serviceError) *trip.CreateTripResponse {
	return &trip.CreateTripResponse{BaseResp: errorBaseResp(err)}
}

func listTripsErrorResponse(err *serviceError) *trip.ListTripsResponse {
	return &trip.ListTripsResponse{BaseResp: errorBaseResp(err)}
}

func latestTripErrorResponse(err *serviceError) *trip.GetLatestTripResponse {
	return &trip.GetLatestTripResponse{BaseResp: errorBaseResp(err)}
}

func detailTripErrorResponse(err *serviceError) *trip.GetTripDetailResponse {
	return &trip.GetTripDetailResponse{BaseResp: errorBaseResp(err)}
}

func deleteTripErrorResponse(err *serviceError) *trip.DeleteTripResponse {
	return &trip.DeleteTripResponse{BaseResp: errorBaseResp(err)}
}

func updateTripErrorResponse(err *serviceError) *trip.UpdateTripResponse {
	return &trip.UpdateTripResponse{BaseResp: errorBaseResp(err)}
}

func toTripDTO(input *tripModel) *trip.TripInfo {
	if input == nil {
		return nil
	}
	return &trip.TripInfo{
		Id:             input.ID,
		UserId:         input.UserID,
		Title:          input.Title,
		Subtitle:       stringPtrIfNotEmpty(input.Subtitle),
		Destination:    stringPtrIfNotEmpty(input.Destination),
		DateRange:      stringPtrIfNotEmpty(input.DateRange),
		DayCount:       int32PtrIfPositive(input.DayCount),
		People:         stringPtrIfNotEmpty(input.People),
		BudgetLevel:    stringPtrIfNotEmpty(input.BudgetLevel),
		SourceQuestion: stringPtrIfNotEmpty(input.SourceQuestion),
		SourceReply:    stringPtrIfNotEmpty(input.SourceReply),
		Summary:        cloneTripSummary(input.Summary),
		Days:           cloneTripDays(input.Days),
		Budget:         cloneTripBudget(input.Budget),
		Alerts:         append([]string(nil), input.Alerts...),
		Saved:          input.Saved,
		CreatedAt:      timePtrString(input.CreatedAt),
		UpdatedAt:      timePtrString(input.UpdatedAt),
	}
}

func toTripDTOList(items []*tripModel) []*trip.TripInfo {
	result := make([]*trip.TripInfo, 0, len(items))
	for _, item := range items {
		result = append(result, toTripDTO(item))
	}
	return result
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

func timePtrString(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.Format(time.RFC3339)
	return &formatted
}

func cloneTripSummary(input *trip.TripSummary) *trip.TripSummary {
	if input == nil {
		return nil
	}
	return &trip.TripSummary{
		Date:   cloneStringPtr(input.Date),
		Days:   cloneStringPtr(input.Days),
		People: cloneStringPtr(input.People),
		Budget: cloneStringPtr(input.Budget),
	}
}

func cloneTripDays(input []*trip.TripDay) []*trip.TripDay {
	result := make([]*trip.TripDay, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripDay{
			Day:     item.Day,
			Title:   cloneStringPtr(item.Title),
			Route:   cloneStringPtr(item.Route),
			Food:    cloneStringPtr(item.Food),
			Hotel:   cloneStringPtr(item.Hotel),
			Tips:    cloneTripTips(item.Tips),
			Weather: cloneStringPtr(item.Weather),
		})
	}
	return result
}

func cloneTripTips(input []*trip.TripTip) []*trip.TripTip {
	result := make([]*trip.TripTip, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripTip{
			Icon:  cloneStringPtr(item.Icon),
			Title: cloneStringPtr(item.Title),
			Text:  cloneStringPtr(item.Text),
		})
	}
	return result
}

func cloneTripBudget(input []*trip.TripBudget) []*trip.TripBudget {
	result := make([]*trip.TripBudget, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		result = append(result, &trip.TripBudget{
			Label:  cloneStringPtr(item.Label),
			Amount: cloneStringPtr(item.Amount),
		})
	}
	return result
}

func cloneStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := *input
	return &value
}
