package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"

	"github.com/google/uuid"
)

type IMerchantItemUsecase interface {
	Create(payload CreateMerchanItemPayload, createdBy string) (entity.MerchantItem, error)
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
