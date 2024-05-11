package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductCheckout struct {
	CheckoutID uuid.UUID `json:"transactionId"`
	CustomerID uuid.UUID `json:"customerId"`
	Paid       int       `json:"paid"`
	Change     int       `json:"change"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ProductCheckoutItem struct {
	CheckoutID uuid.UUID `json:"transactionId"`
	ProductID  uuid.UUID `json:"productId"`
	Amount     int       `json:"amount"`
	Quantity   int       `json:"quantity"`
	CreatedBy  string    `json:"-"`
	CreatedAt  time.Time `json:"-"`
}
