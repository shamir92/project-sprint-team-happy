package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity"
)

type MerchantService struct {
	MerchantZoneRecord        []*entity.MerchantTSPZoneRecord
	MerchantNearestZoneRecord []*entity.MerchantNearestZoneRecord
	PregeneratedMerchants     map[entity.PregeneratedId]*entity.Merchant
	AssignedMerchants         map[entity.PregeneratedId]*entity.Merchant
	merchantSize, totalCount  *int64
	Merchants                 []*entity.Merchant
	SequenceMutex             *sync.Mutex
	merchantShards            []*shard
}

type shard struct {
	merchants []*entity.Merchant
	mutex     *sync.Mutex
	count     int
}

func NewMerchantService() *MerchantService {
	return &MerchantService{
		SequenceMutex: &sync.Mutex{},
		merchantSize:  new(int64),
		totalCount:    new(int64),
	}
}

func (service *MerchantService) GetPregeneratedMerchant() *entity.Merchant {
	for {
		total := atomic.LoadInt64(service.totalCount)
		if total >= int64(*service.merchantSize) {
			return nil
		}

		shardIndex := int(total) % len(service.merchantShards)
		s := service.merchantShards[shardIndex]

		s.mutex.Lock()
		if s.count < len(s.merchants) {
			merchant := s.merchants[s.count]
			s.count = s.count + 1
			atomic.AddInt64(service.totalCount, 1)
			s.mutex.Unlock()
			return merchant
		}
		s.mutex.Unlock()
	}
}

func (service *MerchantService) GetAllMerchantNearestLocations() ([]*entity.MerchantNearestZoneRecord, error) {
	zoneRecords := service.MerchantNearestZoneRecord
	for _, zone := range zoneRecords {
		if !service.checkMerchantNearestRecord(zone) {
			return nil, errors.New("some merchant is still not assigned")
		}
	}
	return service.MerchantNearestZoneRecord, nil
}

func (service *MerchantService) GetMerchantNearestLocations() (*entity.MerchantNearestZoneRecord, error) {
	zoneRecords := service.MerchantNearestZoneRecord
	zone := zoneRecords[rand.IntN(len(zoneRecords))]
	if !service.checkMerchantNearestRecord(zone) {
		return nil, errors.New("some merchant is still not assigned")
	}
	return zone, nil
}

func (service *MerchantService) checkMerchantNearestRecord(zone *entity.MerchantNearestZoneRecord) bool {
	for _, record := range zone.Records {
		for _, merchant := range record.MerchantOrder {
			if merchant.MerchantId == "" {
				return false
			}
		}
	}
	return true
}

func (service *MerchantService) GetAllMerchantRoutes() ([]*entity.MerchantTSPZoneRecord, error) {
	zoneRecords := service.MerchantZoneRecord

	for _, zone := range zoneRecords {
		if !service.checkTspZoneRecord(zone) {
			return nil, errors.New("some merchant is still not assigned")
		}
	}

	return service.MerchantZoneRecord, nil
}

func (service *MerchantService) GetTwoZoneMerchantRoutes() ([]*entity.MerchantTSPZoneRecord, error) {
	zoneRecords := service.MerchantZoneRecord

	randZone1Index := rand.IntN(len(zoneRecords))
	randZone2Index := rand.IntN(len(zoneRecords))
	if randZone1Index == randZone2Index {
		randZone2Index = (randZone2Index + 1) % len(zoneRecords)
	}

	if !service.checkTspZoneRecord(zoneRecords[randZone1Index]) {
		return nil, errors.New("some merchant is still not assigned")
	}
	if !service.checkTspZoneRecord(zoneRecords[randZone2Index]) {
		return nil, errors.New("some merchant is still not assigned")
	}
	zone1 := zoneRecords[randZone1Index]
	zone2 := zoneRecords[randZone2Index]
	return []*entity.MerchantTSPZoneRecord{zone1, zone2}, nil
}

func (service *MerchantService) GetMerchantRoutes() (*entity.MerchantTSPZoneRecord, error) {
	zoneRecords := service.MerchantZoneRecord
	zone := zoneRecords[rand.IntN(len(zoneRecords))]
	if !service.checkTspZoneRecord(zone) {
		return nil, errors.New("some merchant is still not assigned")
	}
	return zone, nil
}

func (service *MerchantService) checkTspZoneRecord(zone *entity.MerchantTSPZoneRecord) bool {
	for _, tspRoute := range zone.GeneratedTSPRoutes {
		for _, merchant := range tspRoute.GeneratedRoutes {
			if merchant.MerchantId == "" {
				return false
			}
		}
	}
	return true

}

// Get all pregenerated merchants
func (service *MerchantService) GetAllPregeneratedMerchants() map[entity.PregeneratedId]*entity.Merchant {
	return service.PregeneratedMerchants
}

func (service *MerchantService) InitMerchantNearestLocations(generateCount int) {
	for a := 0; a < len(service.MerchantZoneRecord); a++ {
		currentZone := service.MerchantZoneRecord[a]
		userLocation := entity.GenerateRandomLocation(entity.LocationPoint{
			Lat:  currentZone.StartZoneRange.Lat,
			Long: currentZone.StartZoneRange.Long,
		}, entity.LocationPoint{
			Lat:  currentZone.EndZoneRange.Lat,
			Long: currentZone.EndZoneRange.Long,
		})
		var records []entity.MerchantNearestRecord
		for i := 0; i < generateCount; i++ {
			merchantPregeneratedIds := currentZone.MerchantPregeneratedId
			merchants := make([]*entity.Merchant, 0)
			distances := make(map[string]float64)

			for _, m := range merchantPregeneratedIds {
				merchant, isExists := service.PregeneratedMerchants[entity.PregeneratedId(m)]
				if isExists {
					distance := entity.CalculateDistance(userLocation, merchant.Location)
					distances[merchant.PregeneratedId] = distance
					merchants = append(merchants, merchant)
				}
			}

			sort.SliceStable(merchants, func(i, j int) bool {
				return distances[merchants[i].PregeneratedId] < distances[merchants[j].PregeneratedId]
			})

			nearestMerchant := make(map[entity.Order]*entity.Merchant)
			for j, merchant := range merchants {
				nearestMerchant[entity.Order(j)] = merchant
			}
			records = append(records, entity.MerchantNearestRecord{
				StartingPoint: userLocation,
				MerchantOrder: nearestMerchant,
			})

			// Log progress for every 100 generated nearest locations
			if (i+1)%100 == 0 {
				log.Printf("Generated %d nearest locations in zone %d\n", i+1, a)
			}
		}
		service.MerchantNearestZoneRecord = append(service.MerchantNearestZoneRecord, &entity.MerchantNearestZoneRecord{
			Records: records,
		})
	}
	log.Println("Init merchant nearest locations")
}

// Initialize pregenerated merchants
func (service *MerchantService) InitPegeneratedTSPMerchants(generateCount int) {
	for a := 0; a < len(service.MerchantZoneRecord); a++ {
		// for each zone, generate random user location
		merchantZone := service.MerchantZoneRecord[a]
		for i := 0; i < generateCount; i++ {
			// generate random user location based on zone bounds
			userLocation := entity.GenerateRandomLocation(entity.LocationPoint{
				Lat: merchantZone.StartZoneRange.Lat, Long: merchantZone.StartZoneRange.Long,
			}, entity.LocationPoint{
				Lat: merchantZone.EndZoneRange.Lat, Long: merchantZone.EndZoneRange.Long,
			})

			// select random 5-10 merchant pregeneratedId from the zone
			selectedMerchant := make([]*entity.Merchant, 0)
			merchantCount := len(merchantZone.MerchantPregeneratedId)
			itemCount := rand.IntN(6) + 5
			rand.Shuffle(merchantCount, func(i, j int) {
				merchantZone.MerchantPregeneratedId[i], merchantZone.MerchantPregeneratedId[j] = merchantZone.MerchantPregeneratedId[j], merchantZone.MerchantPregeneratedId[i]
			})

			// from pregeneratedId, get the merchant
			for j := 0; j < itemCount; j++ {
				preRegeneratedId := entity.PregeneratedId(merchantZone.MerchantPregeneratedId[j])
				merchant, isExists := service.PregeneratedMerchants[preRegeneratedId]
				if isExists {
					selectedMerchant = append(selectedMerchant, merchant)
				}
			}

			// generate TSP route for the selected merchants
			startingIndex := rand.IntN(len(selectedMerchant))
			routes, totalDistance := entity.GenerateTSPMerchantRouteFromStartingPoint(userLocation, startingIndex, selectedMerchant)
			service.MerchantZoneRecord[a].GeneratedTSPRoutes = append(service.MerchantZoneRecord[a].GeneratedTSPRoutes, &entity.MerchantTSPRoute{
				GeneratedRoutes:       routes,
				StartingPoint:         userLocation,
				TotalDistance:         totalDistance,
				TotalDurationInMinute: entity.CalculateTimeInMinute(totalDistance, 40),
				StartingIndex:         startingIndex,
			})

			// Log progress for every 100 generated TSP
			if (i+1)%100 == 0 {
				log.Printf("Generated %d TSP routes in zone %d\n", i+1, a)
			}
		}
	}
	log.Println("Init pregenerated TSP merchants")
}

// Initialize zones with pregenerated merchants
func (service *MerchantService) InitZonesWithPregeneratedMerchants(params entity.MerchantZoneOpts) {
	service.PregeneratedMerchants = make(map[entity.PregeneratedId]*entity.Merchant)
	service.AssignedMerchants = make(map[entity.PregeneratedId]*entity.Merchant)
	merchantShards := make([]*shard, params.NumberOfZones)
	for i := 0; i < params.NumberOfZones; i++ {
		merchantShards[i] = &shard{}
	}
	// init starting point
	var startingPoint entity.LocationPoint
	// create zones contains collection of merchants
	for i := 0; i < params.NumberOfZones; i++ {
		shard := &shard{
			mutex: &sync.Mutex{},
			count: 0,
		}
		// generate flat square location bounds
		endPoint := entity.CalculateAreaSquare(startingPoint, params.Area)

		// create merchants for the zone
		var pregeneratedMerchantIds []string
		for j := 0; j < params.NumberOfMerchantsPerZone; j++ {
			merchants := make(map[string]*entity.Merchant)

			merchant := entity.GenerateRandomMerchant(
				entity.LocationPoint{Lat: startingPoint.Lat, Long: startingPoint.Long},
				entity.LocationPoint{Lat: endPoint.Lat, Long: endPoint.Long},
			)
			merchants[merchant.PregeneratedId] = merchant
			service.PregeneratedMerchants[entity.PregeneratedId(merchant.PregeneratedId)] = merchant
			service.SequenceMutex.Lock()
			service.Merchants = append(service.Merchants, merchant)
			shard.merchants = append(shard.merchants, merchant)
			atomic.AddInt64(service.merchantSize, 1)
			service.SequenceMutex.Unlock()

			pregeneratedMerchantIds = append(pregeneratedMerchantIds, merchant.PregeneratedId)
		}
		merchantShards[i%params.NumberOfZones] = shard

		// append the zone with the merchants into the service
		merchantPregeneratedIds := make([]entity.PregeneratedId, len(pregeneratedMerchantIds))
		for i, id := range pregeneratedMerchantIds {
			merchantPregeneratedIds[i] = entity.PregeneratedId(id)
		}

		service.MerchantZoneRecord = append(service.MerchantZoneRecord, &entity.MerchantTSPZoneRecord{
			StartZoneRange:         entity.LocationPoint{Lat: startingPoint.Lat, Long: startingPoint.Long},
			EndZoneRange:           entity.LocationPoint{Lat: endPoint.Lat, Long: endPoint.Long},
			MerchantPregeneratedId: merchantPregeneratedIds,
		})

		// update the starting point for the next loop
		startingPoint = entity.LocationPoint{Lat: endPoint.Lat, Long: endPoint.Long + entity.KmToDegrees(params.Gap)}
	}
	service.merchantShards = merchantShards
	log.Println("Init zones with pregenerated merchants")
}

// Assign pregenerated merchant to assigned merchant
func (service *MerchantService) AssignPregeneratedMerchant(pregeneratedId, merchantId string) error {
	merchant, isExists := service.PregeneratedMerchants[entity.PregeneratedId(pregeneratedId)]
	if !isExists {
		return fmt.Errorf("pregenerated merchant with id %s does not exist", pregeneratedId)
	}
	merchant.MerchantId = merchantId
	service.AssignedMerchants[entity.PregeneratedId(pregeneratedId)] = merchant
	return nil
}

func (service *MerchantService) ResetAll() {
	service.MerchantZoneRecord = nil
	service.MerchantNearestZoneRecord = nil
	service.PregeneratedMerchants = nil
	service.AssignedMerchants = nil
	service.merchantShards = nil
	service.merchantSize = nil
	service.totalCount = nil
	service.Merchants = nil
	service.SequenceMutex = nil
	log.Println("Reset all")
}
