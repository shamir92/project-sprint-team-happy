package repository

import (
	"database/sql"
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"errors"
)

type productRepository struct {
	db *sql.DB
}
type IProductRepository interface {
	Insert(product entity.Product) (entity.Product, error)
	GetById(id string) (entity.Product, error)
	Update(product entity.Product) error
	Delete(id string) error
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db}
}

func (r *productRepository) Insert(product entity.Product) (entity.Product, error) {
	query := `
		INSERT INTO products(name, sku, category, image_url, notes, price, stock, location, is_available, created_by) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at`

	err := r.db.QueryRow(query, product.Name, product.SKU, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, product.CreatedBy).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (r *productRepository) GetById(id string) (entity.Product, error) {
	var product entity.Product

	query := `
		SELECT 
		id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at, created_by
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.SKU, &product.Category, &product.ImageUrl, &product.Notes, &product.Price, &product.Stock, &product.Location, &product.IsAvailable, &product.CreatedAt, &product.CreatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, commons.CustomError{Message: "product not found", Code: 404}
		}

		return product, err
	}

	return product, nil
}

func (r *productRepository) Update(product entity.Product) error {
	query := `
		UPDATE products SET 
			name = $2,
			sku = $3,
			category = $4,
			image_url = $5,
			notes = $6,
			price = $7,
			stock = $8,
			location = $9,
			is_available = $10,
			updated_at = current_timestamp
		WHERE id = $1 
		AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, product.ID, product.Name, product.SKU, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return commons.CustomError{Message: "product not found", Code: 404}
	}

	return nil
}

func (r *productRepository) Delete(id string) error {
	query := `UPDATE products SET deleted_at = current_timestamp WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.QueryRow(query, id).Err(); err != nil {
		return err
	}

	return nil
}
