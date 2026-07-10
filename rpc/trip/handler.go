package rpctrip

import (
	"context"
	"sync"

	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
)

// TripServiceImpl 实现接口定义中的行程服务。
type TripServiceImpl struct {
	initOnce sync.Once
	service  *tripService
	initErr  error
}

func (s *TripServiceImpl) getService() (*tripService, error) {
	s.initOnce.Do(func() {
		cfg := config.MustGlobal()
		repo, err := mysqlrepo.NewTripRepository(cfg.MySQL)
		if err != nil {
			s.initErr = err
			return
		}
		s.service = newTripService(repo)
	})
	return s.service, s.initErr
}

func (s *TripServiceImpl) CreateTrip(ctx context.Context, req *trip.CreateTripRequest) (*trip.CreateTripResponse, error) {
	service, err := s.getService()
	if err != nil {
		return createTripErrorResponse(errInternal("行程服务初始化失败")), nil
	}
	created, svcErr := service.CreateTrip(ctx, req)
	if svcErr != nil {
		return createTripErrorResponse(svcErr), nil
	}
	return &trip.CreateTripResponse{BaseResp: successBaseResp(), Trip: toTripDTO(created)}, nil
}

func (s *TripServiceImpl) ListTrips(ctx context.Context, req *trip.ListTripsRequest) (*trip.ListTripsResponse, error) {
	service, err := s.getService()
	if err != nil {
		return listTripsErrorResponse(errInternal("行程服务初始化失败")), nil
	}
	items, svcErr := service.ListTrips(ctx, req)
	if svcErr != nil {
		return listTripsErrorResponse(svcErr), nil
	}
	return &trip.ListTripsResponse{BaseResp: successBaseResp(), Trips: toTripDTOList(items)}, nil
}

func (s *TripServiceImpl) GetLatestTrip(ctx context.Context, req *trip.GetLatestTripRequest) (*trip.GetLatestTripResponse, error) {
	if req == nil {
		return latestTripErrorResponse(errParam("请求不能为空")), nil
	}
	service, err := s.getService()
	if err != nil {
		return latestTripErrorResponse(errInternal("行程服务初始化失败")), nil
	}
	item, svcErr := service.GetLatestTrip(ctx, req.UserId)
	if svcErr != nil {
		return latestTripErrorResponse(svcErr), nil
	}
	return &trip.GetLatestTripResponse{BaseResp: successBaseResp(), Trip: toTripDTO(item)}, nil
}

func (s *TripServiceImpl) GetTripDetail(ctx context.Context, req *trip.GetTripDetailRequest) (*trip.GetTripDetailResponse, error) {
	if req == nil {
		return detailTripErrorResponse(errParam("请求不能为空")), nil
	}
	service, err := s.getService()
	if err != nil {
		return detailTripErrorResponse(errInternal("行程服务初始化失败")), nil
	}
	item, svcErr := service.GetTripDetail(ctx, req.UserId, req.TripId)
	if svcErr != nil {
		return detailTripErrorResponse(svcErr), nil
	}
	return &trip.GetTripDetailResponse{BaseResp: successBaseResp(), Trip: toTripDTO(item)}, nil
}

func (s *TripServiceImpl) DeleteTrip(ctx context.Context, req *trip.DeleteTripRequest) (*trip.DeleteTripResponse, error) {
	if req == nil {
		return deleteTripErrorResponse(errParam("请求不能为空")), nil
	}
	service, err := s.getService()
	if err != nil {
		return deleteTripErrorResponse(errInternal("行程服务初始化失败")), nil
	}
	if svcErr := service.DeleteTrip(ctx, req.UserId, req.TripId); svcErr != nil {
		return deleteTripErrorResponse(svcErr), nil
	}
	return &trip.DeleteTripResponse{BaseResp: successBaseResp()}, nil
}
