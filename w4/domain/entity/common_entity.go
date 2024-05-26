package entity

import (
	"math"
)

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

func (st SortType) String() string {
	return string(st)
}

func (st SortType) Valid() bool {
	return st == SortTypeAsc || st == SortTypeDesc
}

const (
	// Earth's radius = 6,371km
	earthRadiusInKm = 6371.0
)

type Location struct {
	Lat float64 `json:"lat" validate:"required,latitude"`
	Lon float64 `json:"long" validate:"required,longitude"`
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

// Distance calculates the Haversine distance between two Location points.
// The distance is returned in kilometers.
func (l Location) Distance(to Location) float64 {
	dLat := degreesToRadians(l.Lat - to.Lat)
	dLon := degreesToRadians(l.Lon - to.Lon)

	// Haversine formula:
	// a = sin²(Δlat/2) + cos(lat1) * cos(lat2) * sin²(Δlon/2)
	a := math.Pow(math.Sin(dLat/2), 2) + math.Pow(math.Sin(dLon/2), 2)*math.Cos(l.Lon)*math.Cos(to.Lon)

	// d = 2R * sin^-1(sqrt(a))
	c := 2 * math.Asin(math.Sqrt(a))
	d := c * earthRadiusInKm

	return d
}
