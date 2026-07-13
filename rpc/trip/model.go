package rpctrip

import (
	"strings"

	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type tripModel = repository.Trip

func defaultTripTitle(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return "我的旅行行程"
	}
	return trimmed
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func optionalInt32(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
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

type updateTripRequest struct {
	Title       *string           `json:"title,omitempty"`
	Subtitle    *string           `json:"subtitle,omitempty"`
	Destination *string           `json:"destination,omitempty"`
	DateRange   *string           `json:"date_range,omitempty"`
	DayCount    *int32            `json:"day_count,omitempty"`
	People      *string           `json:"people,omitempty"`
	BudgetLevel *string           `json:"budget_level,omitempty"`
	Days        []*tripDayPayload `json:"days,omitempty"`
}

func toModelTripDays(input []*tripDayPayload) []*trip.TripDay {
	items := make([]*trip.TripDay, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		items = append(items, &trip.TripDay{
			Day:     item.Day,
			Title:   stringPtrIfNotEmpty(item.Title),
			Route:   stringPtrIfNotEmpty(item.Route),
			Food:    stringPtrIfNotEmpty(item.Food),
			Hotel:   stringPtrIfNotEmpty(item.Hotel),
			Tips:    toModelTripTips(item.Tips),
			Weather: stringPtrIfNotEmpty(item.Weather),
		})
	}
	return items
}

func toModelTripTips(input []tripTipPayload) []*trip.TripTip {
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

func toServiceTripDayPayloads(input []*trip.TripDay) []*tripDayPayload {
	if len(input) == 0 {
		return nil
	}
	items := make([]*tripDayPayload, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		items = append(items, &tripDayPayload{
			Day:     item.Day,
			Title:   item.GetTitle(),
			Route:   item.GetRoute(),
			Food:    item.GetFood(),
			Hotel:   item.GetHotel(),
			Tips:    toServiceTripTips(item.Tips),
			Weather: item.GetWeather(),
		})
	}
	return items
}

func toServiceTripTips(input []*trip.TripTip) []tripTipPayload {
	items := make([]tripTipPayload, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		items = append(items, tripTipPayload{
			Icon:  item.GetIcon(),
			Title: item.GetTitle(),
			Text:  item.GetText(),
		})
	}
	return items
}
