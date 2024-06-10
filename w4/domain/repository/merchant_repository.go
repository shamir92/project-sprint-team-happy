package repository

import (
	"belimang/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mmcloughlin/geohash"
	"github.com/umahmood/haversine"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type IMerchantRepository interface {
	Create(ctx context.Context, merchant entity.Merchant) (entity.Merchant, error)
	Fetch(ctx context.Context, filter entity.MerchantFetchFilter) ([]entity.Merchant, int, error)
	FetchNearby(ctx context.Context, userCoordinate entity.UserCoordinate, neighbors []string, filter entity.MerchantFetchFilter) ([]entity.MerchantWithItem, error)
}

type merchantRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewMerchantRepository(db *sql.DB) *merchantRepository {
	return &merchantRepository{
		db:     db,
		tracer: otel.Tracer("merchant-repository"),
	}
}

func (r *merchantRepository) Create(ctx context.Context, merchant entity.Merchant) (entity.Merchant, error) {
	_, span := r.tracer.Start(ctx, "Create")
	defer span.End()

	merchant.GeoHash = geohash.Encode(merchant.Lat, merchant.Lon)
	sql := "INSERT INTO public.merchants(name, category, image_url, lat, lon, geohash) VALUES ($1, $2, $3, $4, $5,$6) RETURNING id"

	err := r.db.QueryRow(sql, merchant.Name, merchant.Category.String(), merchant.ImageUrl, merchant.Lat, merchant.Lon, merchant.GeoHash).Scan(&merchant.ID)
	if err != nil {
		log.Printf("Failed to create merchant: %v", err.Error())
		return entity.Merchant{}, err
	}

	return merchant, nil
}

func (r *merchantRepository) Fetch(ctx context.Context, filter entity.MerchantFetchFilter) ([]entity.Merchant, int, error) {
	_, span := r.tracer.Start(ctx, "Fetch")
	defer span.End()
	var (
		entities   []entity.Merchant
		conditions []string = make([]string, 0)
		values     []any    = make([]any, 0)
		count               = 0
	)

	sql := "SELECT id, name, category, image_url, lat, lon, geohash FROM public.merchants"
	countQuery := `SELECT count(id) FROM merchants`

	if filter.Name != "" {
		values = append(values, "%"+filter.Name+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(values)))
	}

	if filter.MerchantCategory.String() != "" {
		values = append(values, filter.MerchantCategory.String())
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(values)))
	}

	if filter.Limit <= 0 || filter.Offset < 0 {
		filter.Limit = 5
		filter.Offset = 0
	}

	if len(conditions) > 0 {
		sql += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
		countQuery += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}

	if filter.SortCreatedAt.String() == entity.SortTypeAsc.String() {
		sql += fmt.Sprintf(" ORDER BY created_at ASC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	} else {
		sql += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	}

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		log.Printf("Failed to Fetch merchants: %v", err.Error())
		return entities, count, err
	}

	var _entities []entity.Merchant
	for rows.Next() {
		var entity entity.Merchant

		err = rows.Scan(&entity.ID, &entity.Name, &entity.Category, &entity.ImageUrl, &entity.Lat, &entity.Lon, &entity.GeoHash)
		if err != nil {
			log.Printf("Failed to fetch merchants: %v", err.Error())
			return _entities, count, err
		}

		entities = append(entities, entity)
	}

	err = r.db.QueryRow(countQuery, values...).Scan(&count)

	if err != nil {
		log.Printf("ERROR | MerchantRepository.Fetch() | %v\n", err)
		return _entities, count, err
	}

	return entities, count, nil
}

func (r *merchantRepository) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	_, km := haversine.Distance(haversine.Coord{Lat: lat1, Lon: lon1}, haversine.Coord{Lat: lat2, Lon: lon2})
	return km
}

func (r *merchantRepository) FetchNearby(ctx context.Context, userCoordinate entity.UserCoordinate, neighbors []string, filter entity.MerchantFetchFilter) ([]entity.MerchantWithItem, error) {
	_, span := r.tracer.Start(ctx, "FetchNearby")
	defer span.End()
	var (
		merchantWithItem []entity.MerchantWithItem
		merchants        []entity.Merchant
		conditions       []string = make([]string, 0)
		values           []any    = make([]any, 0)
		geoConditions    string
	)

	sql := "SELECT id, name, category, image_url, lat, lon, geohash FROM public.merchants"
	for i, hash := range neighbors {
		if i == 0 {
			geoConditions = fmt.Sprintf("geohash LIKE '%s%%'", hash)
		} else {
			geoConditions += fmt.Sprintf(" OR geohash LIKE '%s%%'", hash)
		}
	}

	if geoConditions != "" {
		sql += fmt.Sprintf(" WHERE (%s)", geoConditions)
	}

	if filter.Name != "" {
		values = append(values, "%"+filter.Name+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(values)))
	}

	if filter.MerchantCategory.String() != "" {
		values = append(values, filter.MerchantCategory.String())
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(values)))
	}

	// if filter.Limit <= 0 || filter.Offset < 0 {
	// 	filter.Limit = 100
	// 	filter.Offset = 0
	// }

	if len(conditions) > 0 {
		sql += fmt.Sprintf(" AND (%s)", strings.Join(conditions, " AND "))
	}

	if filter.SortCreatedAt.String() == entity.SortTypeAsc.String() {
		// sql += fmt.Sprintf(" ORDER BY created_at ASC")
		sql += " ORDER BY created_at ASC"
	} else {
		// sql += fmt.Sprintf(" ORDER BY created_at DESC", "")
		sql += " ORDER BY created_at DESC"
	}

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		log.Printf("Failed to Fetch merchants: %v", err.Error())
		return nil, err
	}

	for rows.Next() {
		var merchant entity.Merchant
		// var items []entity.MerchantItem

		err = rows.Scan(&merchant.ID, &merchant.Name, &merchant.Category, &merchant.ImageUrl, &merchant.Lat, &merchant.Lon, &merchant.GeoHash)
		if err != nil {
			log.Printf("Failed to fetch merchants: %v", err.Error())
			return nil, err
		}
		merchant.Distance = r.calculateDistance(userCoordinate.Lat, userCoordinate.Lon, merchant.Lat, merchant.Lon)
		merchants = append(merchants, merchant)
	}

	sort.Slice(merchants, func(i, j int) bool {
		return merchants[i].Distance < merchants[j].Distance
	})

	if filter.Limit <= 0 || filter.Offset < 0 {
		filter.Limit = 5
		filter.Offset = 0
	}
	if len(merchants) > filter.Limit {
		merchants = merchants[:filter.Limit]
	}
	log.Println(merchants)
	for _, merchant := range merchants {
		var items []entity.MerchantItem
		sql2 := "SELECT id, name, category, price, image_url, created_at FROM public.merchant_items WHERE merchant_id = $1"
		rows2, err := r.db.Query(sql2, merchant.ID)
		if err != nil {
			log.Printf("Failed to Fetch merchants: %v", err.Error())
			return nil, err
		}
		for rows2.Next() {
			var item entity.MerchantItem
			err = rows2.Scan(&item.ID, &item.Name, &item.Category, &item.Price, &item.ImageUrl, &item.CreatedAt)
			if err != nil {
				log.Printf("Failed to fetch merchants: %v", err.Error())
				return nil, err
			}
			items = append(items, item)
		}

		merchantWithItem = append(merchantWithItem, entity.MerchantWithItem{
			Merchant: merchant,
			Items:    items,
		})
	}

	return merchantWithItem, nil
}
