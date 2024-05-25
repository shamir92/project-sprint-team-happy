package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"strconv"

	"github.com/google/uuid"
)

type IMerchantItemUsecase interface {
	Create(payload CreateMerchanItemPayload, createdBy string) (entity.MerchantItem, error)
	FindItemsByMerchant(query dto.FindMerchantItemPayload) ([]entity.MerchantItem, PaginationMeta, error)
}

type merchantItemUsecase struct {
	itemRepository repository.IMerchantItemRepository
}

func NewMerchanItemUsecase(itemRepository repository.IMerchantItemRepository) *merchantItemUsecase {
	return &merchantItemUsecase{
		itemRepository: itemRepository,
	}
}

type CreateMerchanItemPayload struct {
	Price      int `json:"price" validate:"required,min=1"`
	MerchantID string
	Name       string `json:"name" validate:"required,min=2,max=30"`
	Category   string `json:"productCategory" validate:"required,item_category"`
	ImageUrl   string `json:"imageUrl" validate:"required,urlformat"`
}

type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func (miu *merchantItemUsecase) Create(payload CreateMerchanItemPayload, createdBy string) (entity.MerchantItem, error) {
	var merchantID, err = uuid.Parse(payload.MerchantID)

	if err != nil {
		return entity.MerchantItem{}, helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	isExist, err := miu.itemRepository.CheckMerchantExist(payload.MerchantID)

	if err != nil {
		return entity.MerchantItem{}, err
	}

	if !isExist {
		return entity.MerchantItem{}, helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	newItem, err := miu.itemRepository.Insert(entity.MerchantItem{
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

func (miu *merchantItemUsecase) FindItemsByMerchant(query dto.FindMerchantItemPayload) ([]entity.MerchantItem, PaginationMeta, error) {

	var meta PaginationMeta

	q := query
	if l, err := strconv.Atoi(query.Limit); err == nil {
		meta.Limit = l
	} else {
		// Default Limit
		meta.Limit = 5
		q.Limit = "5"
	}

	if o, err := strconv.Atoi(query.Offset); err == nil {
		meta.Offset = o
	} else {
		// Default Offset
		meta.Offset = 0
		q.Offset = "0"
	}

	rows, count, err := miu.itemRepository.FindAndCount(q)
	meta.Total = count

	return rows, meta, err
}
