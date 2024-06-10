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
	ID        uuid.UUID        `json:"merchantId"`
	Name      string           `json:"name"`
	Category  MerchantCategory `json:"merchantCategory"`
	ImageUrl  string           `json:"imageUrl"`
	Lat       float64          `json:"lat"`
	Lon       float64          `json:"lon"`
	GeoHash   string           `json:"geohash"`
	CreatedAt time.Time        `json:"createdAt"`
}

type MerchantWithItem struct {
	Merchant Merchant       `json:"merchant"`
	Items    []MerchantItem `json:"items"`
}

type MerchantFetchFilter struct {
	ID               string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory"`
	SortCreatedAt    SortType         `json:"createdAt"`
	Limit            int              `json:"limit"`
	Offset           int              `json:"offset"`
}

func (mc Merchant) Location() Location {
	return Location{
		Lat:     mc.Lat,
		Lon:     mc.Lon,
		GeoHash: mc.GeoHash,
	}
}
