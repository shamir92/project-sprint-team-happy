package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/protocol/api/dto"
)

type IMerchantUsecase interface {
	Create(payload MerchantCreatePayload) (dto.MerchantCreateDtoResponse, error)
	Fetch(query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, error)
}

type merchantUsecase struct {
	merchantRepository repository.IMerchantRepository
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository) *merchantUsecase {
	return &merchantUsecase{merchantRepository}
}

type MerchantCreatePayload struct {
	Name     string                  `json:"name" validate:"required,min=2,max=30"`
	Category entity.MerchantCategory `json:"merchantCategory" validate:"required,category=merchant"`
	ImageUrl string                  `json:"imageUrl" validate:"required,urlformat"`
	Location entity.Location         `json:"location" validate:"required"`
}

func (u *merchantUsecase) Create(payload MerchantCreatePayload) (dto.MerchantCreateDtoResponse, error) {
	var entity entity.Merchant

	entity.Name = payload.Name
	entity.Category = payload.Category
	entity.ImageUrl = payload.ImageUrl
	entity.Lat = payload.Location.Lat
	entity.Lon = payload.Location.Lon
	entity, err := u.merchantRepository.Create(entity)
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
	Limit         int                     `query:"limit"`
	Offset        int                     `query:"offset"`
}

func (u *merchantUsecase) Fetch(query MerchantFetchQuery) ([]dto.MerchantFetchDtoResponse, error) {
	var (
		response []dto.MerchantFetchDtoResponse = make([]dto.MerchantFetchDtoResponse, 0)
		filter   entity.MerchantFetchFilter
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

	filter.Limit = query.Limit
	filter.Offset = query.Offset

	merchants, err := u.merchantRepository.Fetch(filter)
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
				Lat: merchant.Lat,
				Lon: merchant.Lon,
			},
		})
	}

	return response, nil
}
