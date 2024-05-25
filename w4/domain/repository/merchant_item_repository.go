package repository

import (
	"belimang/domain/entity"
	"belimang/protocol/api/dto"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	ASC  string = "asc"
	DESC string = "desc"
)

type IMerchantItemRepository interface {
	Insert(entity.MerchantItem) (entity.MerchantItem, error)
	CheckMerchantExist(merchantID string) (bool, error)

	// FindAndCount(query dto.FindMerchantItemPayload) (rows []entity.MerchantItem, count int, err error)
	FindAndCount(query dto.FindMerchantItemPayload) (rows []entity.MerchantItem, count int, err error)
}

type merchantItemRepository struct {
	db *sql.DB
}

func NewMerchanItemRepository(db *sql.DB) *merchantItemRepository {
	return &merchantItemRepository{
		db: db,
	}
}

func (mir *merchantItemRepository) Insert(i entity.MerchantItem) (entity.MerchantItem, error) {
	q := `
		INSERT INTO 
			merchan_items(merchant_id, name, category, price, image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id				
	`

	row := mir.db.QueryRow(q, i.MerchantID, i.Name, i.Category.String(), i.Price, i.ImageUrl, i.CreatedBy)

	var createdItem = i
	if err := row.Scan(&createdItem.ID); err != nil {
		log.Printf("ERROR | MerchanItemRepository.Insert() | %v", err)
	}

	return createdItem, nil
}

func (mir *merchantItemRepository) CheckMerchantExist(merchantID string) (bool, error) {
	var count int
	err := mir.db.QueryRow(`SELECT count(id) FROM merchants WHERE id = $1`, merchantID).Scan(&count)

	return count > 0, err
}

func (mir *merchantItemRepository) FindAndCount(query dto.FindMerchantItemPayload) ([]entity.MerchantItem, int, error) {
	var (
		entities   []entity.MerchantItem = make([]entity.MerchantItem, 0)
		conditions []string              = make([]string, 0)
		values     []any                 = make([]any, 0)
	)

	q := `
		SELECT id, name, category, price, image_url, created_at
		FROM merchant_items mi`

	countQuery := `SELECT count(id) FROM merchant_items mi`

	if err := uuid.Validate(query.MerchantID); err == nil {
		values = append(values, query.MerchantID)
		conditions = append(conditions, fmt.Sprintf("mi.merchant_id = $%d", len(values)))
	}

	if err := uuid.Validate(query.ItemID); err == nil {
		values = append(values, query.ItemID)
		conditions = append(conditions, fmt.Sprintf("mi.id = $%d", len(values)))
	}

	if entity.ValidMerchantItemCategory(query.Category) {
		values = append(values, query.Category)
		conditions = append(conditions, fmt.Sprintf("mi.category = $%d", len(values)))
	}

	if query.Name != "" {
		values = append(values, "%"+query.Name+"%")
		conditions = append(conditions, fmt.Sprintf("mi.name ILIKE $%d", len(values)))
	}

	if len(conditions) > 0 {
		q += fmt.Sprintf("\nWHERE %s", strings.Join(conditions, " AND "))
		countQuery += fmt.Sprintf("\nWHERE %s\n", strings.Join(conditions, " AND "))
	}

	if query.SortCreated == ASC || query.SortCreated == DESC {
		q += fmt.Sprintf(" ORDER BY %s ASC LIMIT %s OFFSET %s", query.SortCreated, query.Limit, query.Offset)
	}

	rows, err := mir.db.Query(q, values...)

	if err != nil {
		log.Printf("ERROR | FindAndCount() | %v", err)
		return entities, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.MerchantItem

		// id, name, category, price, image_url, created_at
		err := rows.Scan(&item.ID, &item.Category, &item.Price, &item.ImageUrl, &item.CreatedAt)

		if err != nil {
			log.Printf("ERROR | FindAndCount() | %v", err)
			return entities, 0, err
		}

		entities = append(entities, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ERROR | FindAndCount() | %v", err)
		return entities, 0, err
	}

	var count int
	err = mir.db.QueryRow(countQuery, values...).Scan(&count)
	if err != nil {
		log.Printf("ERROR | FindAndCount() | %v", err)
		return entities, 0, err
	}

	return entities, count, nil
}
