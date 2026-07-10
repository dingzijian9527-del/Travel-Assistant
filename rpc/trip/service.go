package rpctrip

import (
	"context"
	"strings"
	"time"

	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
)

type tripService struct {
	repo tripRepository
}

func newTripService(repo tripRepository) *tripService {
	return &tripService{repo: repo}
}

func (s *tripService) CreateTrip(ctx context.Context, req *trip.CreateTripRequest) (*tripModel, *serviceError) {
	if req == nil {
		return nil, errParam("请求不能为空")
	}
	if req.UserId <= 0 {
		return nil, errParam("用户信息无效")
	}
	if s == nil || s.repo == nil {
		return nil, errInternal("行程服务不可用")
	}
	now := time.Now()
	input := &tripModel{
		UserID:         req.UserId,
		Title:          defaultTripTitle(optionalString(req.Title)),
		Subtitle:       optionalString(req.Subtitle),
		Destination:    optionalString(req.Destination),
		DateRange:      optionalString(req.DateRange),
		DayCount:       optionalInt32(req.DayCount),
		People:         optionalString(req.People),
		BudgetLevel:    optionalString(req.BudgetLevel),
		SourceQuestion: optionalString(req.SourceQuestion),
		SourceReply:    optionalString(req.SourceReply),
		Summary:        cloneTripSummary(req.Summary),
		Days:           cloneTripDays(req.Days),
		Budget:         cloneTripBudget(req.Budget),
		Alerts:         append([]string(nil), req.Alerts...),
		Saved:          true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	stored, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, errBiz("创建行程失败")
	}
	return stored, nil
}

func (s *tripService) ListTrips(ctx context.Context, req *trip.ListTripsRequest) ([]*tripModel, *serviceError) {
	if req == nil {
		return nil, errParam("请求不能为空")
	}
	if req.UserId <= 0 {
		return nil, errParam("用户信息无效")
	}
	if s == nil || s.repo == nil {
		return nil, errInternal("行程服务不可用")
	}
	limit := int(optionalInt32(req.Limit))
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	items, err := s.repo.ListByUser(ctx, req.UserId, limit)
	if err != nil {
		return nil, errBiz("读取行程列表失败")
	}
	return items, nil
}

func (s *tripService) GetLatestTrip(ctx context.Context, userID int64) (*tripModel, *serviceError) {
	if userID <= 0 {
		return nil, errParam("用户信息无效")
	}
	if s == nil || s.repo == nil {
		return nil, errInternal("行程服务不可用")
	}
	item, ok, err := s.repo.LatestByUser(ctx, userID)
	if err != nil {
		return nil, errBiz("读取最新行程失败")
	}
	if !ok {
		return nil, errNotFound("行程不存在")
	}
	return item, nil
}

func (s *tripService) GetTripDetail(ctx context.Context, userID int64, tripID string) (*tripModel, *serviceError) {
	if userID <= 0 {
		return nil, errParam("用户信息无效")
	}
	if strings.TrimSpace(tripID) == "" {
		return nil, errParam("行程标识无效")
	}
	if s == nil || s.repo == nil {
		return nil, errInternal("行程服务不可用")
	}
	item, ok, err := s.repo.FindByIDForUser(ctx, userID, tripID)
	if err != nil {
		return nil, errBiz("读取行程详情失败")
	}
	if !ok {
		return nil, errNotFound("行程不存在")
	}
	return item, nil
}

func (s *tripService) DeleteTrip(ctx context.Context, userID int64, tripID string) *serviceError {
	if userID <= 0 {
		return errParam("用户信息无效")
	}
	if strings.TrimSpace(tripID) == "" {
		return errParam("行程标识无效")
	}
	if s == nil || s.repo == nil {
		return errInternal("行程服务不可用")
	}
	if !s.repo.DeleteByIDForUser(ctx, userID, tripID) {
		return errNotFound("行程不存在")
	}
	return nil
}
