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
	Category   string `json:"productCategory" validate:"required,category=item"`
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
	var (
		// Default meta
		meta = PaginationMeta{
			Limit:  5,
			Offset: 0,
			Total:  0,
		}
		q = query
	)

	if err := miu.validateMerchantIsExist(query.MerchantID); err != nil {
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

	rows, count, err := miu.itemRepository.FindAndCount(q)
	meta.Total = count

	return rows, meta, err
}

func (miu *merchantItemUsecase) validateMerchantIsExist(id string) error {
	var _, err = uuid.Parse(id)

	if err != nil {
		return helper.CustomError{
			Message: "merchant not found",
			Code:    404,
		}
	}

	isExist, err := miu.itemRepository.CheckMerchantExist(id)

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
