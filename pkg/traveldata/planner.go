package traveldata

import (
	"context"
	"strings"
)

type Planner struct {
	mapProvider     MapProvider
	weatherProvider WeatherProvider
}

func NewPlanner(mapProvider MapProvider, weatherProvider WeatherProvider) *Planner {
	return &Planner{mapProvider: mapProvider, weatherProvider: weatherProvider}
}

func (p *Planner) BuildContext(ctx context.Context, req Request) (Result, error) {
	destination := strings.TrimSpace(req.Destination)
	if destination == "" {
		destination = InferDestination(req.Message)
	}
	days := normalizePositive(req.Days, normalizePositive(InferDays(req.Message), 1))
	people := normalizePositive(req.People, normalizePositive(InferPeople(req.Message), 1))
	budget := req.Budget
	if budget <= 0 {
		budget = InferBudget(req.Message)
	}
	result := Result{
		Destination: destination,
		DateRange:   strings.TrimSpace(req.DateRange),
		Days:        days,
		People:      people,
		BudgetItems: SplitBudget(BudgetRequest{Total: budget, Days: days, People: people}),
	}
	if destination == "" {
		result.Alerts = append(result.Alerts, "未识别到明确目的地，已按通用旅行预算和热门方向给出参考")
		return result, nil
	}
	city := City{Name: destination}
	if p != nil && p.mapProvider != nil {
		resolved, err := p.mapProvider.ResolveCity(ctx, destination)
		if err == nil && strings.TrimSpace(resolved.Name) != "" {
			city = resolved
			result.Destination = resolved.Name
		}
	}
	if p != nil && p.weatherProvider != nil {
		weather, err := p.weatherProvider.Forecast(ctx, city, result.Days)
		if err == nil {
			result.Weather = weather
		}
	}
	if p != nil && p.mapProvider != nil {
		result.Places = append(result.Places, p.search(ctx, city, "景点", 2)...)
		result.Places = append(result.Places, p.search(ctx, city, "美食", 2)...)
		result.Places = append(result.Places, p.search(ctx, city, "住宿", 1)...)
		if len(result.Places) >= 2 {
			route, err := p.mapProvider.Route(ctx, result.Places[0], result.Places[1])
			if err == nil && route.From != "" && route.To != "" {
				result.Routes = append(result.Routes, route)
			}
		}
	}
	if len(result.Weather) == 0 {
		result.Alerts = append(result.Alerts, "天气数据暂不可用，出发前请以实时天气为准")
	}
	if len(result.Places) == 0 {
		result.Alerts = append(result.Alerts, "地点数据暂不可用，热门景点和餐饮请以平台实时信息为准")
	}
	return result, nil
}

func (p *Planner) search(ctx context.Context, city City, keyword string, limit int) []Place {
	places, err := p.mapProvider.SearchPlaces(ctx, city, keyword, limit)
	if err != nil {
		return nil
	}
	if limit > 0 && len(places) > limit {
		return places[:limit]
	}
	return places
}

func normalizePositive(value int, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
