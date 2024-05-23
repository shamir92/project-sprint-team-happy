package entity

import (
	"time"

	"github.com/google/uuid"
)

type ItemCategory string

const (
	Beverage   ItemCategory = "Beverage"
	Food       ItemCategory = "Food"
	Snack      ItemCategory = "Snack"
	Condiments ItemCategory = "Condiments"
	Additions  ItemCategory = "Additions"
)

func (pc ItemCategory) String() string {
	return string(pc)
}

type MerchantItem struct {
	ID        uuid.UUID
	Name      string
	Category  ItemCategory
	ImageUrl  string
	Price     int
	CreatedAt time.Time
}
