package dto

import (
	"time"

	"github.com/google/uuid"
)

type FindMerchantItemPayload struct {
	ItemID      string
	Limit       string
	Offset      string
	Name        string
	Category    string
	SortCreated string
	MerchantID  string
}

type MerchanItemDto struct {
	ID        uuid.UUID `json:"itemId"`
	Name      string    `json:"name"`
	Category  string    `json:"productCategory"`
	ImageUrl  string    `json:"imageUrl"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type FindMerchanItemMetaDto struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
