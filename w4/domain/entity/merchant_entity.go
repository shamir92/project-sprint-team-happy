package entity

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

var merchantCategories []MerchantCategory = []MerchantCategory{
	SmallRestaurant,
	MediumRestaurant,
	LargeRestaurant,
	MerchandiseRestaurant,
	BoothKiosk,
	ConvenienceStore,
}

func (mc MerchantCategory) String() string {
	return string(mc)
}

func (mc MerchantCategory) Valid() bool {
	return slices.Index(merchantCategories, MerchantCategory(mc)) != -1
}

type Merchant struct {
	ID        uuid.UUID
	Name      string
	Category  MerchantCategory
	ImageUrl  string
	Lat       float64
	Lon       float64
	CreatedAt time.Time
}

type MerchantFetchFilter struct {
	ID               string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory"`
	SortCreatedAt    SortType         `json:"createdAt"`
	Limit            int              `json:"limit"`
	Offset           int              `json:"offset"`
}
