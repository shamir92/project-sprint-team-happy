package dto

import (
	"belimang/domain/entity"
	"time"

	"github.com/google/uuid"
)

type OrderEstimateResponseDto struct {
	TotalPrice                     int       `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int       `json:"estimatedDeliveryTimeInMinutes"`
	OrderId                        uuid.UUID `json:"calculatedEstimateId"`
}

type PlaceOrderRequestDto struct {
	OrderId string `json:"calculatedEstimateId"`
}

type PlaceOrderResponseDto struct {
	OrderId string `json:"orderId"`
}

type GetOrderSearchParams struct {
	MerchantID string `json:"merchant_id"`
	Name       string `json:"name"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Category   string `json:"category"`
}

type OrderItemDto struct {
	ItemId          uuid.UUID           `json:"itemId"`
	Name            string              `json:"name"`
	ProductCategory entity.ItemCategory `json:"productCategory"`
	Price           int                 `json:"price"`
	Quantity        int                 `json:"quantity"`
	ImageUrl        string              `json:"imageUrl"`
	CreatedAt       time.Time           `json:"createdAt"`
}

type OrderDetailDto struct {
	Merchant MerchantFetchDtoResponse `json:"merchant"`
	Items    []OrderItemDto           `json:"items"`
}

type GetOrderResponseDto struct {
	OrderID string           `json:"orderId"`
	Orders  []OrderDetailDto `json:"orders"`
}
