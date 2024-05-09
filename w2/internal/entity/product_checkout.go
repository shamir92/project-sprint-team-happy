package entity

import "time"

type ProductCheckout struct {
	ID         string `json:"userId"`
	CustomerID string `json:"name"`
	Paid       int    `json:"phoneNumber"`
	Change     int    `json:"change"`
	CreatedBy  string `json:"-"`
	CreatedAt  time.Time
	// UpdatedAt
}

type ProductCheckoutItem struct {
}
