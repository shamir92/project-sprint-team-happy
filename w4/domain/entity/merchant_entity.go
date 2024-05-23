package entity

import (
	"time"

	"github.com/google/uuid"
)

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "Small Restaurant"
	MediumRestaurant      MerchantCategory = "Medium Restaurant"
	LargeRestaurant       MerchantCategory = "Large Restaurant"
	MerchandiseRestaurant MerchantCategory = "Merchandise Restaurant"
	BoothKiosk            MerchantCategory = "Booth Kiosk"
	ConvenienceStore      MerchantCategory = "Convenience Store"
)

func (mc MerchantCategory) String() string {
	return string(mc)
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
