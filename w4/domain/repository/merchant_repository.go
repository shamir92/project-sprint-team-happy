package repository

import (
	"belimang/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mmcloughlin/geohash"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type IMerchantRepository interface {
	Create(ctx context.Context, merchant entity.Merchant) (entity.Merchant, error)
	Fetch(ctx context.Context, filter entity.MerchantFetchFilter) ([]entity.Merchant, error)
	FetchNearby(ctx context.Context, geohash []string, filter entity.MerchantFetchFilter) ([]entity.Merchant, error)
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

func (r *merchantRepository) Fetch(ctx context.Context, filter entity.MerchantFetchFilter) ([]entity.Merchant, error) {
	_, span := r.tracer.Start(ctx, "Fetch")
	defer span.End()
	var (
		entities   []entity.Merchant
		conditions []string = make([]string, 0)
		values     []any    = make([]any, 0)
	)

	sql := "SELECT id, name, category, image_url, lat, lon, geohash FROM public.merchants"

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
	}

	if filter.SortCreatedAt.String() == entity.SortTypeAsc.String() {
		sql += fmt.Sprintf(" ORDER BY created_at ASC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	} else {
		sql += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	}

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		log.Printf("Failed to Fetch merchants: %v", err.Error())
		return entities, err
	}

	var _entities []entity.Merchant
	for rows.Next() {
		var entity entity.Merchant

		err = rows.Scan(&entity.ID, &entity.Name, &entity.Category, &entity.ImageUrl, &entity.Lat, &entity.Lon, &entity.GeoHash)
		if err != nil {
			log.Printf("Failed to fetch merchants: %v", err.Error())
			return _entities, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *merchantRepository) FetchNearby(ctx context.Context, neighbors []string, filter entity.MerchantFetchFilter) ([]entity.Merchant, error) {
	_, span := r.tracer.Start(ctx, "Fetch")
	defer span.End()
	var (
		entities   []entity.Merchant
		conditions []string = make([]string, 0)
		values     []any    = make([]any, 0)
	)

	var sql string
	for i, hash := range neighbors {
		if i == 0 {
			sql = fmt.Sprintf("SELECT id, name, category, image_url, lat, lon, geohash FROM public.merchants WHERE geohash LIKE '%s%%'", hash)
		} else {
			sql += fmt.Sprintf(" OR geohash LIKE '%s%%'", hash)
		}
	}
	// conditions = append(conditions, fmt.Sprintf("category = $%d", len(values)))

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
		sql += fmt.Sprintf(" %s", strings.Join(conditions, " AND "))
	}

	if filter.SortCreatedAt.String() == entity.SortTypeAsc.String() {
		sql += fmt.Sprintf(" ORDER BY created_at ASC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	} else {
		sql += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", filter.Limit, filter.Offset)
	}
	log.Println(sql)
	log.Println(values)
	rows, err := r.db.Query(sql, values...)
	if err != nil {
		log.Printf("Failed to Fetch merchants: %v", err.Error())
		return entities, err
	}

	var _entities []entity.Merchant
	for rows.Next() {
		var entity entity.Merchant

		err = rows.Scan(&entity.ID, &entity.Name, &entity.Category, &entity.ImageUrl, &entity.Lat, &entity.Lon, &entity.GeoHash)
		if err != nil {
			log.Printf("Failed to fetch merchants: %v", err.Error())
			return _entities, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}
