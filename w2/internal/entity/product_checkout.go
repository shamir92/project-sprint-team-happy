package entity

import (
	"time"

	"github.com/google/uuid"
)

// type ProductCheckout struct {
// 	ID         string `json:"userId"`
// 	CustomerID string `json:"name"`
// 	Paid       int    `json:"phoneNumber"`
// 	Change     int    `json:"change"`
// 	CreatedBy  string `json:"-"`
// 	CreatedAt  time.Time
// 	// UpdatedAt
// }

type ProductCheckoutItem struct {
	CheckoutID uuid.UUID `json:"name"`
	ProductID  uuid.UUID `json:"phoneNumber"`
	Amount     int       `json:"amount"`
	Quantity   int       `json:"quantity"`
	CreatedBy  string    `json:"-"`
	CreatedAt  time.Time
}
