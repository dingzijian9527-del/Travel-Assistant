package traveldata

import (
	"context"
	"fmt"
	"strings"
)

type Request struct {
	Message     string
	Destination string
	DateRange   string
	Days        int
	People      int
	Budget      int
	Preferences []string
}

type City struct {
	Name     string
	Adcode   string
	Location string
}

type WeatherDay struct {
	Date    string
	Text    string
	TempMin int
	TempMax int
	Wind    string
}

type Place struct {
	Name     string
	Address  string
	Category string
	Location string
}

type Route struct {
	From            string
	To              string
	DistanceMeters  int
	DurationMinutes int
}

type BudgetItem struct {
	Label  string
	Amount int
	Note   string
}

type Result struct {
	Destination string
	DateRange   string
	Days        int
	People      int
	Weather     []WeatherDay
	Places      []Place
	Routes      []Route
	BudgetItems []BudgetItem
	Alerts      []string
}

type MapProvider interface {
	ResolveCity(ctx context.Context, destination string) (City, error)
	SearchPlaces(ctx context.Context, city City, keyword string, limit int) ([]Place, error)
	Route(ctx context.Context, from Place, to Place) (Route, error)
}

type WeatherProvider interface {
	Forecast(ctx context.Context, city City, days int) ([]WeatherDay, error)
}

func (r Result) FormatForPrompt() string {
	if r.Destination == "" && len(r.Weather) == 0 && len(r.Places) == 0 && len(r.BudgetItems) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("旅行实时参考：\n")
	if r.Destination != "" {
		fmt.Fprintf(&b, "- 目的地：%s\n", r.Destination)
	}
	if r.DateRange != "" {
		fmt.Fprintf(&b, "- 出行时间：%s\n", r.DateRange)
	}
	if r.Days > 0 {
		fmt.Fprintf(&b, "- 游玩天数：%d天\n", r.Days)
	}
	if r.People > 0 {
		fmt.Fprintf(&b, "- 出行人数：%d人\n", r.People)
	}
	if len(r.Weather) > 0 {
		b.WriteString("- 天气参考：")
		for i, item := range r.Weather {
			if i > 0 {
				b.WriteString("；")
			}
			fmt.Fprintf(&b, "%s %s %d-%d℃ %s", item.Date, item.Text, item.TempMin, item.TempMax, item.Wind)
		}
		b.WriteString("\n")
	}
	if len(r.Places) > 0 {
		b.WriteString("- 地点推荐：")
		for i, item := range r.Places {
			if i > 0 {
				b.WriteString("；")
			}
			fmt.Fprintf(&b, "%s（%s，%s）", item.Name, item.Category, item.Address)
		}
		b.WriteString("\n")
	}
	if len(r.Routes) > 0 {
		b.WriteString("- 路线参考：")
		for i, item := range r.Routes {
			if i > 0 {
				b.WriteString("；")
			}
			fmt.Fprintf(&b, "%s到%s约%d公里，约%d分钟", item.From, item.To, item.DistanceMeters/1000, item.DurationMinutes)
		}
		b.WriteString("\n")
	}
	if len(r.BudgetItems) > 0 {
		b.WriteString("- 预算拆分：")
		for i, item := range r.BudgetItems {
			if i > 0 {
				b.WriteString("；")
			}
			fmt.Fprintf(&b, "%s约%d元", item.Label, item.Amount)
		}
		b.WriteString("\n")
	}
	if len(r.Alerts) > 0 {
		b.WriteString("- 出行提醒：")
		b.WriteString(strings.Join(r.Alerts, "；"))
		b.WriteString("\n")
	}
	b.WriteString("请基于以上参考生成回答；如果实时数据不足，说明以官方和平台实时信息为准。\n")
	return b.String()
}
