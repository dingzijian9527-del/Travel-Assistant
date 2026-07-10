package traveldata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type QWeatherClient struct {
	key     string
	baseURL string
	client  *http.Client
}

func NewQWeatherClient(key string, timeout time.Duration) *QWeatherClient {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &QWeatherClient{
		key:     key,
		baseURL: "https://devapi.qweather.com",
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *QWeatherClient) Forecast(ctx context.Context, city City, days int) ([]WeatherDay, error) {
	if c == nil || c.key == "" || city.Location == "" {
		return nil, errors.New("和风天气输入不完整")
	}
	values := url.Values{}
	values.Set("key", c.key)
	values.Set("location", city.Location)
	var payload struct {
		Code  string `json:"code"`
		Daily []struct {
			FxDate  string `json:"fxDate"`
			TextDay string `json:"textDay"`
			TempMin string `json:"tempMin"`
			TempMax string `json:"tempMax"`
			WindDir string `json:"windDirDay"`
		} `json:"daily"`
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v7/weather/3d?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("和风天气接口状态异常")
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Code != "200" {
		return nil, errors.New("和风天气返回异常")
	}
	limit := normalizePositive(days, 3)
	items := make([]WeatherDay, 0, len(payload.Daily))
	for i, item := range payload.Daily {
		if i >= limit {
			break
		}
		items = append(items, WeatherDay{
			Date:    item.FxDate,
			Text:    item.TextDay,
			TempMin: atoi(item.TempMin),
			TempMax: atoi(item.TempMax),
			Wind:    item.WindDir,
		})
	}
	return items, nil
}
