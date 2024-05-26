package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderState string

const (
	Estimated OrderState = "estimated"
	Ordered   OrderState = "ordered"
)

func (os OrderState) String() string {
	return string(os)
}

type Order struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	UserLat               float64
	UserLon               float64
	TotalPrice            int
	EstimatedDeliveryTime int
	State                 OrderState
	CreatedAt             time.Time
}
