package entity

import (
	"math"
	"math/rand/v2"
)

type LocationZone struct {
	StartPoint     LocationPoint
	EndPoint       LocationPoint
	LocationPoints []LocationPoint
}

type LocationPoint struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Calculate the distance in degrees for a given distance in kilometers
func KmToDegrees(distanceKm float64) float64 {
	const earthRadiusKm = 6371.0
	return distanceKm / earthRadiusKm * (180 / math.Pi)
}

// CalculateBottomRightPoint calculates the bottom-right point of a square on the Earth's surface
// starting from a given top-left point and an area in square kilometers
func CalculateAreaSquare(startPoint LocationPoint, areaKm2 float64) LocationPoint {
	// Calculate the side length of the square in kilometers
	sideLengthKm := math.Sqrt(areaKm2)

	// Convert side length to degrees
	sideLengthDeg := KmToDegrees(sideLengthKm)

	// Calculate the bottom-right corner of the square
	bottomRight := LocationPoint{
		Lat:  startPoint.Lat - sideLengthDeg,
		Long: startPoint.Long + sideLengthDeg,
	}

	return bottomRight
}

func GenerateRandomLocation(startingPoint, endPoint LocationPoint) LocationPoint {
	latitude := rand.Float64()*(endPoint.Lat-startingPoint.Lat) + startingPoint.Lat
	longitude := rand.Float64()*(endPoint.Long-startingPoint.Long) + startingPoint.Long

	return LocationPoint{
		Lat:  latitude,
		Long: longitude,
	}
}

func CalculateTimeInMinute(distanceInKm float64, speedInKmHour int) int {
	return int(distanceInKm / float64(speedInKmHour) * 60)
}

func CalculateDistance(p1, p2 LocationPoint) float64 {
	const R = 6371 // Radius of Earth in kilometers

	lat1 := p1.Lat * math.Pi / 180
	long1 := p1.Long * math.Pi / 180
	lat2 := p2.Lat * math.Pi / 180
	long2 := p2.Long * math.Pi / 180

	dlat := lat2 - lat1
	dlong := long2 - long1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlong/2)*math.Sin(dlong/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
