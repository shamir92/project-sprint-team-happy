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
	return slices.Index(itemCategories, ItemCategory(pc)) != -1
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
	ID         uuid.UUID    `json:"itemId"`
	Name       string       `json:"name"`
	Category   ItemCategory `json:"category"`
	ImageUrl   string       `json:"imageUrl"`
	Price      int          `json:"price"`
	CreatedAt  time.Time    `json:"createdAt"`
	CreatedBy  string       `json:"createdBy"`
	MerchantID uuid.UUID    `json:"merchantId"`
	merchant   Merchant     `json:"-"`
}

func (m *MerchantItem) SetMerchant(merchant Merchant) {
	m.merchant = merchant
}

func (m MerchantItem) Merchant() Merchant {
	return m.merchant
}
