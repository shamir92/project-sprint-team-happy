package service

import (
	"eniqlostore/internal/entity"
	"eniqlostore/internal/repository"
	"strconv"
)

type ProductCheckoutService struct {
	checkoutRepository repository.IProductCheckoutRepository
}

func NewProductCheckoutService(repo repository.IProductCheckoutRepository) *ProductCheckoutService {
	return &ProductCheckoutService{
		checkoutRepository: repo,
	}
}

type FindCheckoutHistoryRequest struct {
	Limit         string `json:"limit"`
	Offset        string `json:"offset"`
	SortCreatedAt string `json:"createdAt"`
	CustomerID    string `json:"customerId"`
}

type FindCheckoutHistoryResponse struct {
	entity.ProductCheckout
	Items []entity.ProductCheckoutItem `json:"productDetails"`
}

func (pcs *ProductCheckoutService) GetHistory(payload FindCheckoutHistoryRequest) ([]FindCheckoutHistoryResponse, error) {
	var options []entity.FindCheckoutHistoryBuilder

	var limit, offset = 5, 0

	if l, err := strconv.Atoi(payload.Limit); err == nil && l > 0 {
		limit = l
	}

	if o, err := strconv.Atoi(payload.Offset); err == nil && o > 0 {
		offset = o
	}

	options = append(options, entity.FindCheckoutWithLimitAndOffset(offset, limit))

	if payload.CustomerID != "" {
		options = append(options, entity.FindCheckoutWithCustomerId(payload.CustomerID))
	}

	// Default newest
	if payload.SortCreatedAt == "" {
		options = append(options, entity.FindCheckoutWithSortByCreatedAt(entity.DESC))
	}

	if payload.SortCreatedAt == entity.ASC.String() {
		options = append(options, entity.FindCheckoutWithSortByCreatedAt(entity.ASC))
	}

	histories, err := pcs.checkoutRepository.Find(options...)

	if err != nil {
		return nil, err
	}

	var checkoutIds []string

	for _, history := range histories {
		checkoutIds = append(checkoutIds, history.CheckoutID.String())
	}

	items, err := pcs.checkoutRepository.FindItemByCheckoutIds(checkoutIds)

	if err != nil {
		return nil, err
	}

	var itemsByCheckoutID = make(map[string][]entity.ProductCheckoutItem, len(checkoutIds))

	for _, item := range items {
		checkoutID := item.CheckoutID.String()

		checkoutItems, ok := itemsByCheckoutID[checkoutID]

		if !ok {
			itemsByCheckoutID[checkoutID] = []entity.ProductCheckoutItem{item}
		} else {
			checkoutItems = append(checkoutItems, item)
			itemsByCheckoutID[checkoutID] = checkoutItems
		}
	}

	var res = make([]FindCheckoutHistoryResponse, 0, len(checkoutIds))

	for _, history := range histories {
		res = append(res, FindCheckoutHistoryResponse{
			ProductCheckout: history,
			Items:           itemsByCheckoutID[history.CheckoutID.String()],
		})
	}

	return res, nil
}
