package traveldata

import "github.com/dingzijian9527-del/Travel-Assistant/pkg/config"

func NewPlannerFromConfig(cfg config.TravelDataConfig) *Planner {
	if !cfg.Enabled {
		return nil
	}
	return NewPlanner(
		NewAmapClient(cfg.AmapKey, cfg.Timeout),
		NewQWeatherClient(cfg.QWeatherKey, cfg.Timeout),
	)
}
