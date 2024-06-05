package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/protocol/api/dto"
	"context"
	"strconv"

	"github.com/mmcloughlin/geohash"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type IMerchantUsecase interface {
	Create(ctx context.Context, payload MerchantCreatePayload) (dto.MerchantCreateDtoResponse, error)
	Fetch(ctx context.Context, query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, entity.PaginationMeta, error)
	FetchNearby(ctx context.Context, userCoordinate UserCoordinate, query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, error)
}

type merchantUsecase struct {
	merchantRepository repository.IMerchantRepository
	tracer             trace.Tracer
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository) *merchantUsecase {
	return &merchantUsecase{
		merchantRepository: merchantRepository,
		tracer:             otel.Tracer("merchant-usecase"),
	}
}

type MerchantCreatePayload struct {
	Name     string                  `json:"name" validate:"required,min=2,max=30"`
	Category entity.MerchantCategory `json:"merchantCategory" validate:"required,category=merchant"`
	ImageUrl string                  `json:"imageUrl" validate:"required,urlformat"`
	Location entity.Location         `json:"location" validate:"required"`
}

func (u *merchantUsecase) Create(ctx context.Context, payload MerchantCreatePayload) (dto.MerchantCreateDtoResponse, error) {
	_, span := u.tracer.Start(ctx, "Create")
	defer span.End()
	var entity entity.Merchant

	entity.Name = payload.Name
	entity.Category = payload.Category
	entity.ImageUrl = payload.ImageUrl
	entity.Lat = payload.Location.Lat
	entity.Lon = payload.Location.Lon
	entity, err := u.merchantRepository.Create(ctx, entity)
	if err != nil {
		return dto.MerchantCreateDtoResponse{}, err
	}

	return dto.MerchantCreateDtoResponse{ID: entity.ID.String()}, nil
}

type MerchantFetchQuery struct {
	ID            string                  `query:"merchantId"`
	Name          string                  `query:"name"`
	Category      entity.MerchantCategory `query:"merchantCategory"`
	ImageUrl      string                  `query:"imageUrl"`
	Location      entity.Location         `query:"location"`
	SortCreatedAt entity.SortType         `query:"createdAt"`
	Limit         string                  `query:"limit"`
	Offset        string                  `query:"offset"`
}

func (u *merchantUsecase) Fetch(ctx context.Context, query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, entity.PaginationMeta, error) {
	_, span := u.tracer.Start(ctx, "Fetch")
	defer span.End()
	var (
		response []dto.MerchantFetchDtoResponse = make([]dto.MerchantFetchDtoResponse, 0)
		filter   entity.MerchantFetchFilter
		meta     = entity.PaginationMeta{
			Limit:  5, // default limit
			Offset: 0, // default offset
			Total:  0,
		}
	)

	if query.ID != "" {
		filter.ID = query.ID
	}

	if query.Name != "" {
		filter.Name = query.Name
	}

	if query.Category.String() != "" {
		if query.Category.Valid() {
			filter.MerchantCategory = query.Category
		} else {
			return response, meta, nil
		}
	}

	if query.SortCreatedAt.Valid() {
		filter.SortCreatedAt = query.SortCreatedAt
	}

	if l, err := strconv.Atoi(query.Limit); err == nil {
		meta.Limit = l
	}

	if o, err := strconv.Atoi(query.Offset); err == nil {
		meta.Offset = o
	}

	filter.Limit = meta.Limit
	filter.Offset = meta.Offset

	merchants, count, err := u.merchantRepository.Fetch(ctx, filter)
	if err != nil {
		return response, meta, err
	}
	meta.Total = count

	for _, merchant := range merchants {
		response = append(response, dto.MerchantFetchDtoResponse{
			ID:       merchant.ID,
			Name:     merchant.Name,
			Category: merchant.Category,
			ImageUrl: merchant.ImageUrl,
			Location: entity.Location{
				Lat:     merchant.Lat,
				Lon:     merchant.Lon,
				GeoHash: merchant.GeoHash,
			},
		})
	}

	return response, meta, nil
}

type UserCoordinate struct {
	Lat float64 `json:"lat" validate:"required,latitude"`
	Lon float64 `json:"long" validate:"required,longitude"`
}

func (u *merchantUsecase) FetchNearby(ctx context.Context, userCoordinate UserCoordinate, query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, error) {
	_, span := u.tracer.Start(ctx, "FetchNearby")
	defer span.End()
	targetGeoHash := geohash.EncodeWithPrecision(userCoordinate.Lat, userCoordinate.Lon, 6)
	neighbors := geohash.Neighbors(targetGeoHash)
	neighbors = append(neighbors, targetGeoHash)

	var (
		response []dto.MerchantFetchDtoResponse = make([]dto.MerchantFetchDtoResponse, 0)
		filter   entity.MerchantFetchFilter
		meta     = entity.PaginationMeta{
			Limit:  5,
			Offset: 0,
		}
	)

	if query.ID != "" {
		filter.ID = query.ID
	}

	if query.Name != "" {
		filter.Name = query.Name
	}

	if query.Category.String() != "" {
		if query.Category.Valid() {
			filter.MerchantCategory = query.Category
		} else {
			return response, nil
		}
	}

	if query.SortCreatedAt.Valid() {
		filter.SortCreatedAt = query.SortCreatedAt
	}

	if l, err := strconv.Atoi(query.Limit); err == nil {
		meta.Limit = l
	}

	if o, err := strconv.Atoi(query.Offset); err == nil {
		meta.Offset = o
	}

	filter.Limit = meta.Limit
	filter.Offset = meta.Offset

	merchants, err := u.merchantRepository.FetchNearby(ctx, neighbors, filter)
	if err != nil {
		return response, err
	}

	for _, merchant := range merchants {
		response = append(response, dto.MerchantFetchDtoResponse{
			ID:       merchant.ID,
			Name:     merchant.Name,
			Category: merchant.Category,
			ImageUrl: merchant.ImageUrl,
			Location: entity.Location{
				Lat:     merchant.Lat,
				Lon:     merchant.Lon,
				GeoHash: merchant.GeoHash,
			},
		})
	}

	return response, nil
}
