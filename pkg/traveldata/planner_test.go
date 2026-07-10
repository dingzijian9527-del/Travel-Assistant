package traveldata

import (
	"context"
	"strings"
	"testing"
)

func TestPlannerBuildsTravelContextFromProviders(t *testing.T) {
	planner := NewPlanner(fakeMapProvider{}, fakeWeatherProvider{})
	ctx := context.Background()

	result, err := planner.BuildContext(ctx, Request{
		Message:     "帮我规划成都三天行程，预算3000元，喜欢美食",
		Destination: "成都",
		Days:        3,
		People:      2,
		Budget:      3000,
	})
	if err != nil {
		t.Fatalf("build context returned error: %v", err)
	}
	if result.Destination != "成都" {
		t.Fatalf("unexpected destination: %s", result.Destination)
	}
	if len(result.Weather) == 0 {
		t.Fatal("expected weather data")
	}
	if len(result.Places) == 0 {
		t.Fatal("expected place data")
	}
	if len(result.Routes) == 0 {
		t.Fatal("expected route data")
	}
	if len(result.BudgetItems) == 0 {
		t.Fatal("expected budget items")
	}
	formatted := result.FormatForPrompt()
	for _, expected := range []string{"旅行实时参考", "成都", "多云", "宽窄巷子", "预算拆分", "路线参考"} {
		if !strings.Contains(formatted, expected) {
			t.Fatalf("formatted context missing %s: %s", expected, formatted)
		}
	}
}

type fakeMapProvider struct{}

func (fakeMapProvider) ResolveCity(ctx context.Context, destination string) (City, error) {
	return City{Name: destination, Location: "104.0668,30.5728", Adcode: "510100"}, nil
}

func (fakeMapProvider) SearchPlaces(ctx context.Context, city City, keyword string, limit int) ([]Place, error) {
	return []Place{
		{Name: "宽窄巷子", Address: "青羊区", Category: keyword, Location: "104.05,30.67"},
		{Name: "春熙路", Address: "锦江区", Category: keyword, Location: "104.08,30.65"},
	}, nil
}

func (fakeMapProvider) Route(ctx context.Context, from Place, to Place) (Route, error) {
	return Route{From: from.Name, To: to.Name, DistanceMeters: 3200, DurationMinutes: 18}, nil
}

type fakeWeatherProvider struct{}

func (fakeWeatherProvider) Forecast(ctx context.Context, city City, days int) ([]WeatherDay, error) {
	return []WeatherDay{{Date: "2026-07-06", Text: "多云", TempMin: 24, TempMax: 32, Wind: "微风"}}, nil
}
