package entity

import (
	"math"
)

type Merchant struct {
	MerchantId     string        `json:"merchantId"`
	PregeneratedId string        `json:"pregeneratedId"`
	Location       LocationPoint `json:"location"`
}

type MerchantTSPRoute struct {
	GeneratedRoutes       map[Order]*Merchant
	StartingPoint         LocationPoint
	TotalDistance         float64
	TotalDurationInMinute int
	StartingIndex         int
}

type MerchantTSPZoneRecord struct {
	StartZoneRange         LocationPoint
	EndZoneRange           LocationPoint
	MerchantPregeneratedId []PregeneratedId
	GeneratedTSPRoutes     []*MerchantTSPRoute
}

type MerchantZoneOpts struct {
	Area                     float64
	Gap                      float64
	NumberOfZones            int
	NumberOfMerchantsPerZone int
}

type MerchantNearestRecord struct {
	StartingPoint LocationPoint
	MerchantOrder map[Order]*Merchant
}

type MerchantNearestZoneRecord struct {
	Records []MerchantNearestRecord
}

type PregeneratedId string
type MerchantId string
type Order int

// Generate random merchant with pregenerated id
func GenerateRandomMerchant(startingPoint, endPoint LocationPoint) *Merchant {
	res := GenerateRandomLocation(startingPoint, endPoint)
	return &Merchant{
		PregeneratedId: GeneratePregeneratedId(),
		Location:       res,
	}
}

// Generate TSP route for merchants from starting point
func GenerateTSPMerchantRouteFromStartingPoint(
	startingPoint LocationPoint, startingPointIndex int, merchants []*Merchant) (map[Order]*Merchant, float64) {

	route := make(map[Order]*Merchant, len(merchants))
	visited := make([]bool, len(merchants))
	currentMerchant := merchants[startingPointIndex]
	totalDistance := CalculateDistance(startingPoint, currentMerchant.Location) // Calculate initial distance

	visited[startingPointIndex] = true
	route[0] = currentMerchant

	for i := 1; i < len(merchants); i++ {
		nearestIndex := -1
		minDistance := math.MaxFloat64

		for j := 0; j < len(merchants); j++ {
			if visited[j] {
				continue
			}

			distance := CalculateDistance(currentMerchant.Location, merchants[j].Location)
			if distance < minDistance {
				minDistance = distance
				nearestIndex = j
			}
		}

		if nearestIndex != -1 {
			route[Order(i)] = merchants[nearestIndex]
			visited[nearestIndex] = true
			currentMerchant = merchants[nearestIndex]
			totalDistance += minDistance
		}
	}

	totalDistance += CalculateDistance(currentMerchant.Location, startingPoint)

	return route, totalDistance
}
