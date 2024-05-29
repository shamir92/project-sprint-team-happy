package entity

import "github.com/google/uuid"

type OrderItem struct {
	OrderID    uuid.UUID
	ItemID     uuid.UUID
	MerchantID uuid.UUID
	Price      int
	Quantity   int
	Amount     int
	Item       *MerchantItem
}

func (o *OrderItem) SetItem(item *MerchantItem) {
	o.Item = item
}
