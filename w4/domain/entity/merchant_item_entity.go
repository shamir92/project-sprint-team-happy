package entity

import (
	"slices"
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

var itemCategories []ItemCategory = []ItemCategory{Beverage, Food, Snack, Condiments, Additions}

func (pc ItemCategory) String() string {
	return string(pc)
}

func (pc ItemCategory) Valid() bool {
	return slices.Index(itemCategories, ItemCategory(pc)) == -1
}

func ValidMerchantItemCategory(category string) bool {
	var categories = []ItemCategory{
		Beverage,
		Food,
		Snack,
		Condiments,
		Additions,
	}

	return slices.IndexFunc(categories, func(ic ItemCategory) bool {
		return ic.String() == category
	}) != -1
}

type MerchantItem struct {
	ID         uuid.UUID
	Name       string
	Category   ItemCategory
	ImageUrl   string
	Price      int
	CreatedAt  time.Time
	CreatedBy  string
	MerchantID uuid.UUID
}
