package traveldata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AmapClient struct {
	key     string
	baseURL string
	client  *http.Client
}

func NewAmapClient(key string, timeout time.Duration) *AmapClient {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &AmapClient{
		key:     key,
		baseURL: "https://restapi.amap.com",
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *AmapClient) ResolveCity(ctx context.Context, destination string) (City, error) {
	if c == nil || c.key == "" {
		return City{Name: destination}, errors.New("高德密钥为空")
	}
	values := url.Values{}
	values.Set("key", c.key)
	values.Set("address", destination)
	var payload struct {
		Status   string `json:"status"`
		Geocodes []struct {
			FormattedAddress string `json:"formatted_address"`
			City             string `json:"city"`
			Adcode           string `json:"adcode"`
			Location         string `json:"location"`
		} `json:"geocodes"`
	}
	if err := c.get(ctx, "/v3/geocode/geo", values, &payload); err != nil {
		return City{Name: destination}, err
	}
	if payload.Status != "1" || len(payload.Geocodes) == 0 {
		return City{Name: destination}, errors.New("高德地理编码失败")
	}
	item := payload.Geocodes[0]
	name := item.City
	if name == "" {
		name = item.FormattedAddress
	}
	if name == "" {
		name = destination
	}
	return City{Name: name, Adcode: item.Adcode, Location: item.Location}, nil
}

func (c *AmapClient) SearchPlaces(ctx context.Context, city City, keyword string, limit int) ([]Place, error) {
	if c == nil || c.key == "" {
		return nil, errors.New("高德密钥为空")
	}
	values := url.Values{}
	values.Set("key", c.key)
	values.Set("keywords", keyword)
	values.Set("city", firstNonEmpty(city.Adcode, city.Name))
	values.Set("offset", strconv.Itoa(normalizePositive(limit, 5)))
	values.Set("page", "1")
	values.Set("extensions", "base")
	var payload struct {
		Status string `json:"status"`
		POIs   []struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Address  any    `json:"address"`
			Location string `json:"location"`
		} `json:"pois"`
	}
	if err := c.get(ctx, "/v3/place/text", values, &payload); err != nil {
		return nil, err
	}
	if payload.Status != "1" {
		return nil, errors.New("高德地点搜索失败")
	}
	places := make([]Place, 0, len(payload.POIs))
	for _, item := range payload.POIs {
		places = append(places, Place{
			Name:     item.Name,
			Address:  stringify(item.Address),
			Category: keyword,
			Location: item.Location,
		})
	}
	return places, nil
}

func (c *AmapClient) Route(ctx context.Context, from Place, to Place) (Route, error) {
	if c == nil || c.key == "" || from.Location == "" || to.Location == "" {
		return Route{}, errors.New("路线输入不完整")
	}
	values := url.Values{}
	values.Set("key", c.key)
	values.Set("origin", from.Location)
	values.Set("destination", to.Location)
	values.Set("strategy", "0")
	var payload struct {
		Status string `json:"status"`
		Route  struct {
			Paths []struct {
				Distance string `json:"distance"`
				Duration string `json:"duration"`
			} `json:"paths"`
		} `json:"route"`
	}
	if err := c.get(ctx, "/v3/direction/driving", values, &payload); err != nil {
		return Route{}, err
	}
	if payload.Status != "1" || len(payload.Route.Paths) == 0 {
		return Route{}, errors.New("高德路线规划失败")
	}
	path := payload.Route.Paths[0]
	distance := atoi(path.Distance)
	durationSeconds := atoi(path.Duration)
	return Route{From: from.Name, To: to.Name, DistanceMeters: distance, DurationMinutes: durationSeconds / 60}, nil
}

func (c *AmapClient) get(ctx context.Context, path string, values url.Values, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+values.Encode(), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("高德接口状态异常")
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func stringify(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		if len(typed) == 0 {
			return ""
		}
		return stringify(typed[0])
	default:
		return ""
	}
}
