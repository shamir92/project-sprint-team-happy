package repository

import (
	"database/sql"
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"errors"
	"fmt"
)

type productRepository struct {
	db *sql.DB
}
type IProductRepository interface {
	Insert(product entity.Product) (entity.Product, error)
	GetById(id string) (entity.Product, error)
	Update(product entity.Product) error
	Delete(id string) error
	Find(...entity.FindProductOptionBuilder) ([]entity.Product, error)
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

func (r *productRepository) Find(opts ...entity.FindProductOptionBuilder) ([]entity.Product, error) {
	options := &entity.FindProductOption{}
	for _, o := range opts {
		o(options)
	}

	query := `
		SELECT
			id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at
		FROM
			products p
		WHERE 1=1
	`

	values := []interface{}{
		options.Limit,
		options.Offset,
	}

	if options.ID != "" {
		values = append(values, options.ID)
		query += fmt.Sprintf(" AND p.id = $%d", len(values))
	}

	if options.IsAvailable != nil {
		values = append(values, *options.IsAvailable)
		query += fmt.Sprintf(" AND p.is_available = $%d", len(values))
	}

	if options.InStock != nil {
		if *options.InStock {
			query += " AND p.stock > 0"
		} else {
			query += " AND p.stock <= 0"
		}
	}

	if options.Name != "" {
		values = append(values, fmt.Sprintf("%%%s%%", options.Name))
		query += fmt.Sprintf(" AND p.name ILIKE $%d", len(values))
	}

	if options.Category != "" {
		values = append(values, options.Category)
		query += fmt.Sprintf(" AND p.category = $%d", len(values))
	}

	if options.SKU != "" {
		values = append(values, options.SKU)
		query += fmt.Sprintf(" AND p.sku = $%d", len(values))
	}

	sorting := map[string]entity.SortType{} // key: column, val: desc | asc

	if options.SortCreatedAt.String() != "" {
		sorting["p.created_at"] = options.SortCreatedAt
	}

	if options.SortPrice.String() != "" {
		sorting["p.price"] = options.SortPrice
	}

	if len(sorting) != 0 {
		var sortQuery = "ORDER BY"
		var count = 0
		var last = len(sorting)
		for key, val := range sorting {
			count++

			if count == last {
				sortQuery += fmt.Sprintf(" %s %s", key, val)
			} else {
				sortQuery += fmt.Sprintf(" %s %s,", key, val)
			}
		}

		fmt.Println(sortQuery)

		query = fmt.Sprintf("%s\n%s", query, sortQuery)
	}

	query = fmt.Sprintf("%s\nLIMIT $1 OFFSET $2", query)
	rows, err := r.db.Query(query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []entity.Product = []entity.Product{}

	for rows.Next() {
		var product entity.Product

		err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Category, &product.ImageUrl, &product.Notes, &product.Price, &product.Stock, &product.Location, &product.IsAvailable, &product.CreatedAt)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
