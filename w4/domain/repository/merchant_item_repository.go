package repository

import (
	"belimang/domain/entity"
	"database/sql"
	"log"
)

type IMerchantItemRepository interface {
	Insert(entity.MerchantItem) (entity.MerchantItem, error)
	CheckMerchantExist(merchantID string) (bool, error)
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
