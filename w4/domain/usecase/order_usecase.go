package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"log"
	"math"

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

	if startPointCounter == 0 || startPointCounter > 1 {
		return emptyOrder, helper.CustomError{
			Message: "only 1 merchant starting point allowed",
			Code:    400,
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
		merchantLocations     []entity.Location
		merchantStartingPoint entity.Location
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
			merchantStartingPoint = merchant.Location()
		} else {
			merchantLocations = append(merchantLocations, merchant.Location())
		}
	}

	userLocation := entity.Location{
		Lat: payload.UserLocation.Lat,
		Lon: payload.UserLocation.Lon,
	}

	estimatedDeliveryTime, err := calculateEstimateOrderDeliveryTime(append([]entity.Location{merchantStartingPoint}, merchantLocations...), userLocation)

	if err != nil {
		return emptyOrder, err
	}

	newOrder := entity.Order{
		UserLat:               payload.UserLocation.Lat,
		UserLon:               payload.UserLocation.Lon,
		TotalPrice:            totalPrice,
		EstimatedDeliveryTime: int(estimatedDeliveryTime), // TODO: calculate estimated delivery time
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

// The first element in merchant locations is assumed as starting point
func calculateEstimateOrderDeliveryTime(merchantLocations []entity.Location, userLocation entity.Location) (float64, error) {
	const SPEED_PER_HOUR = 40    // 40Km
	const MAX_DISTANCE_IN_KM = 3 // 3Km

	var currentLocation = merchantLocations[0]
	var totalDistanceInKm float64

	// Calculate distance start from first merchant til last merchant
	var i = 1
	for len(merchantLocations) > 0 && i < len(merchantLocations) {
		loc := merchantLocations[i]
		var distance float64 = currentLocation.Distance(loc)

		// Check if the distance between the current merchant and the user exceeds the maximum allowed distance
		if currentLocation.Distance(userLocation) > MAX_DISTANCE_IN_KM {
			return 0, helper.CustomError{
				Code:    400,
				Message: "merchant's destination too far from user's location",
			}
		}

		totalDistanceInKm += distance
		currentLocation = loc
		i += 1
	}

	// Calcuate distance from last merchant's location to user's location
	totalDistanceInKm += currentLocation.Distance(userLocation)

	var deliveryTimeInMinutes = (totalDistanceInKm / SPEED_PER_HOUR) * 60
	return math.Round(deliveryTimeInMinutes), nil
}
