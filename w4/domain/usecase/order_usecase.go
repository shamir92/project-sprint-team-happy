package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	errOrderNotFound = errors.New("order not found")
)

type IOrderUsecase interface {
	MakeOrderEstimate(ctx context.Context, payload MakeOrderEstimatePayload, userId string) (entity.Order, error)
	PlaceOrder(ctx context.Context, orderId string, userId string) (entity.Order, error)
	GetOrders(ctx context.Context, params dto.GetOrderSearchParams, userID string) ([]dto.GetOrderResponseDto, error)
}

type orderUsecase struct {
	orderRepository        repository.IOrderRepository
	merchantItemRepository repository.IMerchantItemRepository
	tracer                 trace.Tracer
}

func NewOrderUsecase(orderRepository repository.IOrderRepository, itemRepository repository.IMerchantItemRepository) *orderUsecase {
	return &orderUsecase{
		orderRepository:        orderRepository,
		merchantItemRepository: itemRepository,
		tracer:                 otel.Tracer("order-usecase"),
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

func (u *orderUsecase) MakeOrderEstimate(ctx context.Context, payload MakeOrderEstimatePayload, userId string) (entity.Order, error) {
	_, span := u.tracer.Start(ctx, "MakeOrderEstimate")
	defer span.End()
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

	merchantItems, err := u.merchantItemRepository.FindByItemIds(ctx, itemIds)

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

	estimatedDeliveryTime, err := u.calculateEstimateOrderDeliveryTime(ctx, append([]entity.Location{merchantStartingPoint}, merchantLocations...), userLocation)

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

	createdOrder, err := u.orderRepository.InsertEstimateOrder(ctx, repository.InsertEstimateOrderPayload{
		Order:      newOrder,
		OrderItems: orderItems,
	})

	if err != nil {
		log.Printf("ERROR | MakeOrderEstimate() | %v", err)
		return emptyOrder, err
	}

	return createdOrder, nil
}

func (u *orderUsecase) PlaceOrder(ctx context.Context, orderId string, userId string) (entity.Order, error) {
	_, span := u.tracer.Start(ctx, "PlaceOrder")
	defer span.End()
	if err := uuid.Validate(orderId); err != nil {
		return entity.Order{}, helper.CustomError{
			Message: errOrderNotFound.Error(),
			Code:    404,
		}
	}

	order, err := u.orderRepository.FindOrderByID(ctx, orderId)

	if err != nil {
		return entity.Order{}, err
	}

	if order.UserID.String() != userId {
		log.Printf("hey")
		return entity.Order{}, helper.CustomError{
			Message: errOrderNotFound.Error(),
			Code:    404,
		}
	}

	order.ChangeStateToOrdered()

	if err := u.orderRepository.Update(ctx, order); err != nil {
		return order, err
	}

	return order, nil
}

func (u *orderUsecase) GetOrders(ctx context.Context, params dto.GetOrderSearchParams, userID string) ([]dto.GetOrderResponseDto, error) {
	_, span := u.tracer.Start(ctx, "PlaceOrder")
	defer span.End()
	user, _ := uuid.Parse(userID)

	orders, err := u.orderRepository.FindByUser(ctx, params, user)

	if err != nil {
		var errMsg = fmt.Errorf("ERROR | OrderUsecase.GetOrders() | error to find orders: %v", err)
		log.Println(errMsg.Error())
		return []dto.GetOrderResponseDto{}, errMsg
	}

	var ordersOut []dto.GetOrderResponseDto
	for _, order := range orders {
		var orderResp dto.GetOrderResponseDto
		orderResp.OrderID = order.ID.String()
		orderResp.Orders = make([]dto.OrderDetailDto, 0)
		var orderItemsByMerchant = make(map[uuid.UUID]dto.OrderDetailDto, 0)

		/**
			{
				[merchantId]: { merchant: Merchant, Items: []OrderItem }
			}
		**/
		for _, orderItem := range order.Items {
			orderDetail, ok := orderItemsByMerchant[orderItem.Item.MerchantID]
			var merchant = orderItem.Item.Merchant()
			if !ok {
				orderDetail.Merchant = dto.MerchantFetchDtoResponse{
					ID:        merchant.ID,
					Name:      merchant.Name,
					ImageUrl:  merchant.ImageUrl,
					Location:  merchant.Location(),
					Category:  merchant.Category,
					CreatedAt: merchant.CreatedAt,
				}
			}

			orderDetail.Items = append(orderDetail.Items, dto.OrderItemDto{
				ItemId:          orderItem.ItemID,
				Name:            orderItem.Item.Name,
				ProductCategory: orderItem.Item.Category,
				Price:           orderItem.Price,
				Quantity:        orderItem.Quantity,
				ImageUrl:        orderItem.Item.ImageUrl,
				CreatedAt:       orderItem.Item.CreatedAt,
			})

			orderItemsByMerchant[orderItem.Item.MerchantID] = orderDetail
		}

		for _, orderDetail := range orderItemsByMerchant {
			orderResp.Orders = append(orderResp.Orders, orderDetail)
		}

		ordersOut = append(ordersOut, orderResp)
	}
	return ordersOut, nil
}

// The first element in merchant locations is assumed as starting point
func (u *orderUsecase) calculateEstimateOrderDeliveryTime(ctx context.Context, merchantLocations []entity.Location, userLocation entity.Location) (float64, error) {
	_, span := u.tracer.Start(ctx, "MakeOrderEstimate")
	defer span.End()
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
