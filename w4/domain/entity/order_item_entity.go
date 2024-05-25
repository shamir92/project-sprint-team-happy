package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	ItemID     uuid.UUID
	MerchantID uuid.UUID
	Price      int
	Quantity   int
	Amount     int
}
