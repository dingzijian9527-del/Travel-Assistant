package rpctrip

import (
	"context"
	"testing"

	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
)

func TestTripServiceCreatesQueriesAndDeletesUserScopedTrips(t *testing.T) {
	ctx := context.Background()
	service := newTripService(newTripRepo())

	created, svcErr := service.CreateTrip(ctx, &trip.CreateTripRequest{
		UserId:         42,
		Title:          stringPtr("成都 · 3天"),
		Subtitle:       stringPtr("太古里 / 熊猫基地 / 宽窄巷子"),
		Destination:    stringPtr("成都"),
		DateRange:      stringPtr("7月3日到7月5日"),
		DayCount:       int32Ptr(3),
		People:         stringPtr("2人"),
		BudgetLevel:    stringPtr("舒适型"),
		SourceQuestion: stringPtr("请根据以下信息为我生成专属旅行行程"),
		SourceReply:    stringPtr("第1天：太古里。第2天：熊猫基地。第3天：宽窄巷子。"),
		Summary: &trip.TripSummary{
			Date:   stringPtr("7月3日到7月5日"),
			Days:   stringPtr("3天"),
			People: stringPtr("2人"),
			Budget: stringPtr("舒适型"),
		},
		Days: []*trip.TripDay{
			{
				Day:     1,
				Title:   stringPtr("抵达成都"),
				Route:   stringPtr("太古里 -> 人民公园"),
				Food:    stringPtr("火锅"),
				Hotel:   stringPtr("春熙路附近"),
				Weather: stringPtr("多云"),
				Tips: []*trip.TripTip{
					{Icon: stringPtr("行"), Title: stringPtr("交通建议"), Text: stringPtr("地铁优先")},
				},
			},
		},
		Budget: []*trip.TripBudget{
			{Label: stringPtr("住宿"), Amount: stringPtr("¥1200")},
			{Label: stringPtr("餐饮"), Amount: stringPtr("¥800")},
		},
		Alerts: []string{"熊猫基地建议提前预约"},
	})
	if svcErr != nil {
		t.Fatalf("create trip returned error: %v", svcErr)
	}
	if created.ID == "" || created.UserID != 42 {
		t.Fatalf("unexpected created trip: %+v", created)
	}

	latest, svcErr := service.GetLatestTrip(ctx, 42)
	if svcErr != nil {
		t.Fatalf("latest trip returned error: %v", svcErr)
	}
	if latest.ID != created.ID || latest.Title != "成都 · 3天" {
		t.Fatalf("latest trip mismatch: %+v", latest)
	}

	trips, svcErr := service.ListTrips(ctx, &trip.ListTripsRequest{UserId: 42})
	if svcErr != nil {
		t.Fatalf("list trips returned error: %v", svcErr)
	}
	if len(trips) != 1 || trips[0].ID != created.ID {
		t.Fatalf("unexpected trip list: %+v", trips)
	}

	detail, svcErr := service.GetTripDetail(ctx, 42, created.ID)
	if svcErr != nil {
		t.Fatalf("trip detail returned error: %v", svcErr)
	}
	if len(detail.Days) != 1 || len(detail.Budget) != 2 || len(detail.Alerts) != 1 {
		t.Fatalf("trip detail did not preserve plan data: %+v", detail)
	}
	if len(detail.Days[0].Tips) != 1 || detail.Days[0].GetWeather() != "多云" {
		t.Fatalf("trip detail did not preserve display fields: %+v", detail.Days[0])
	}

	if _, svcErr := service.GetTripDetail(ctx, 7, created.ID); svcErr == nil {
		t.Fatal("another user should not see this trip")
	}
	if svcErr := service.DeleteTrip(ctx, 7, created.ID); svcErr == nil {
		t.Fatal("another user should not delete this trip")
	}
	if svcErr := service.DeleteTrip(ctx, 42, created.ID); svcErr != nil {
		t.Fatalf("delete trip returned error: %v", svcErr)
	}
	if _, svcErr := service.GetTripDetail(ctx, 42, created.ID); svcErr == nil {
		t.Fatal("deleted trip should not be visible")
	}
}

func stringPtr(value string) *string {
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}
