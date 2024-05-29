package dto

import (
	"belimang/domain/entity"
	"time"

	"github.com/google/uuid"
)

type MerchantFetchDtoResponse struct {
	ID        uuid.UUID               `json:"merchantId"`
	Name      string                  `json:"name"`
	Category  entity.MerchantCategory `json:"merchantCategory"`
	ImageUrl  string                  `json:"imageUrl"`
	Location  entity.Location         `json:"location"`
	CreatedAt time.Time               `json:"createdAt,omitempty"`
}

type MerchantCreateDtoResponse struct {
	ID string `json:"merchantId"`
}
