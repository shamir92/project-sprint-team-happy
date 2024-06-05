package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"
	"strconv"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type IMerchantItemUsecase interface {
	Create(ctx context.Context, payload CreateMerchanItemPayload, createdBy string) (entity.MerchantItem, error)
	FindItemsByMerchant(ctx context.Context, query dto.FindMerchantItemPayload) ([]entity.MerchantItem, entity.PaginationMeta, error)
}

type merchantItemUsecase struct {
	itemRepository repository.IMerchantItemRepository
	tracer         trace.Tracer
}

func NewMerchanItemUsecase(itemRepository repository.IMerchantItemRepository) *merchantItemUsecase {
	return &merchantItemUsecase{
		itemRepository: itemRepository,
		tracer:         otel.Tracer("merchant-item-usecase"),
	}
}

type CreateMerchanItemPayload struct {
	Price      int `json:"price" validate:"required,min=1"`
	MerchantID string
	Name       string `json:"name" validate:"required,min=2,max=30"`
	Category   string `json:"productCategory" validate:"required,category=item"`
	ImageUrl   string `json:"imageUrl" validate:"required,urlformat"`
}

func (u *merchantItemUsecase) Create(ctx context.Context, payload CreateMerchanItemPayload, createdBy string) (entity.MerchantItem, error) {
	_, span := u.tracer.Start(ctx, "Create")
	defer span.End()
	var merchantID, err = uuid.Parse(payload.MerchantID)

	if err != nil {
		return entity.MerchantItem{}, helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	isExist, err := u.itemRepository.CheckMerchantExist(ctx, payload.MerchantID)

	if err != nil {
		return entity.MerchantItem{}, err
	}

	if !isExist {
		return entity.MerchantItem{}, helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	newItem, err := u.itemRepository.Insert(ctx, entity.MerchantItem{
		Name:       payload.Name,
		Category:   entity.ItemCategory(payload.Category),
		ImageUrl:   payload.ImageUrl,
		Price:      payload.Price,
		MerchantID: merchantID,
		CreatedBy:  createdBy,
	})

	if err != nil {
		return entity.MerchantItem{}, err
	}

	return newItem, err
}

func (u *merchantItemUsecase) FindItemsByMerchant(ctx context.Context, query dto.FindMerchantItemPayload) ([]entity.MerchantItem, entity.PaginationMeta, error) {
	_, span := u.tracer.Start(ctx, "FindItemsByMerchant")
	defer span.End()
	var (
		// Default meta
		meta = entity.PaginationMeta{
			Limit:  5,
			Offset: 0,
			Total:  0,
		}
		q = query
	)

	if err := u.validateMerchantIsExist(ctx, query.MerchantID); err != nil {
		return nil, meta, err
	}

	if l, err := strconv.Atoi(query.Limit); err == nil {
		meta.Limit = l
	} else {
		q.Limit = strconv.Itoa(meta.Limit)
	}

	if o, err := strconv.Atoi(query.Offset); err == nil {
		meta.Offset = o
	} else {
		q.Offset = strconv.Itoa(meta.Offset)
	}

	if query.SortCreated == entity.SortTypeAsc.String() || query.SortCreated == entity.SortTypeDesc.String() {
		q.SortCreated = query.SortCreated
	} else {
		q.SortCreated = entity.SortTypeDesc.String() // default
	}

	rows, count, err := u.itemRepository.FindAndCount(ctx, q)
	meta.Total = count

	return rows, meta, err
}

func (u *merchantItemUsecase) validateMerchantIsExist(ctx context.Context, id string) error {
	_, span := u.tracer.Start(ctx, "validateMerchantIsExist")
	defer span.End()
	var _, err = uuid.Parse(id)

	if err != nil {
		return helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	isExist, err := u.itemRepository.CheckMerchantExist(ctx, id)

	if err != nil {
		return err
	}

	if !isExist {
		return helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	return nil
}
