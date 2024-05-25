package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type IOrderUsecase interface {
	MakeOrderEstimate(payload MakeOrderEstimatePayload, userId string) (entity.Order, error)
}

type orderUsecase struct {
	orderRepository        repository.IOrderRepository
	merchantItemRepository repository.IMerchantItemRepository
}

func NewOrderUsecase(orderRepository repository.IOrderRepository, itemRepository repository.IMerchantItemRepository) *orderUsecase {
	return &orderUsecase{
		orderRepository:        orderRepository,
		merchantItemRepository: itemRepository,
	}
}

type OrderEstimateItem struct {
	Quantity int    `json:"quantity" validate:"required,min=1"`
	ItemID   string `json:"itemId" validate:"required"`
}

type MakeOrderEstimateMerchant struct {
	IsStartingPoint bool                `json:"isStartingPoint" validate:"required"`
	MerchantID      string              `json:"merchantId" validate:"required"`
	Items           []OrderEstimateItem `json:"items"`
}

type MakeOrderEstimatePayload struct {
	UserLocation entity.Location             `json:"userLocation" validate:"required"`
	Orders       []MakeOrderEstimateMerchant `json:"orders" validate:"required"`
}

func (o *orderUsecase) MakeOrderEstimate(payload MakeOrderEstimatePayload, userId string) (entity.Order, error) {
	var (
		itemIds    = make([]string, 0)
		emptyOrder entity.Order // negative return
	)

	userID, err := uuid.Parse(userId)
	if err != nil {
		log.Printf("ERROR | MakeOrderEstimate() | failed to parse user's id: %v", err)
		return emptyOrder, err
	}

	var startPointCounter int
	for _, order := range payload.Orders {
		if order.IsStartingPoint {
			startPointCounter += 1
		}

		if startPointCounter > 1 {
			return emptyOrder, helper.CustomError{
				Message: "only 1 starting point allowed",
				Code:    400,
			}
		}

		if err := uuid.Validate(order.MerchantID); err != nil {
			return emptyOrder, helper.CustomError{
				Message: "merchat not found",
				Code:    404,
			}
		}

		for _, item := range order.Items {
			if err := uuid.Validate(item.ItemID); err != nil {
				return emptyOrder, helper.CustomError{
					Message: "item not found",
					Code:    404,
				}
			}

			itemIds = append(itemIds, item.ItemID)
		}
	}

	merchantItems, err := o.merchantItemRepository.FindByItemIds(itemIds)

	if err != nil {
		log.Printf("ERROR | MakeOrderEstimate() | %v", err)
		return emptyOrder, err
	}

	itemsByID := map[string]entity.MerchantItem{}

	for _, item := range merchantItems {
		itemsByID[item.ID.String()] = item
	}

	var (
		orderItems            []entity.OrderItem = make([]entity.OrderItem, 0)
		merchants             []entity.Merchant
		merchantStartingPoint entity.Merchant
		totalPrice            int = 0
	)
	for _, order := range payload.Orders {
		var merchant entity.Merchant
		for _, item := range order.Items {
			merchantItem, ok := itemsByID[item.ItemID]
			if !ok {
				return emptyOrder, helper.CustomError{
					Message: "item not found",
					Code:    404,
				}
			}

			if merchantItem.MerchantID.String() != order.MerchantID {
				return emptyOrder, helper.CustomError{
					Message: "merchat not found",
					Code:    404,
				}
			}

			merchant = merchantItem.Merchant()

			totalPrice += merchantItem.Price * item.Quantity
			orderItems = append(orderItems, entity.OrderItem{
				Quantity:   item.Quantity,
				ItemID:     merchantItem.ID,
				MerchantID: merchantItem.MerchantID,
				Price:      merchantItem.Price,
				Amount:     merchantItem.Price * item.Quantity,
			})
		}

		if order.IsStartingPoint {
			merchantStartingPoint = merchant
		} else {
			merchants = append(merchants, merchant)
		}
	}

	fmt.Printf("Merchant Starting Point - lat: %f  lon: %f\n", merchantStartingPoint.Lat, merchantStartingPoint.Lon)
	fmt.Println(merchants)
	newOrder := entity.Order{
		UserLat:               payload.UserLocation.Lat,
		UserLon:               payload.UserLocation.Lon,
		TotalPrice:            totalPrice,
		EstimatedDeliveryTime: 0, // TODO: calculate estimated delivery time
		State:                 entity.Estimated,
		UserID:                userID,
	}

	createdOrder, err := o.orderRepository.InsertEstimateOrder(repository.InsertEstimateOrderPayload{
		Order:      newOrder,
		OrderItems: orderItems,
	})

	if err != nil {
		log.Printf("ERROR | MakeOrderEstimate() | %v", err)
		return emptyOrder, err
	}

	return createdOrder, nil
}
