package repository

import (
	"database/sql"
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	GetByIds(ids []uuid.UUID) ([]entity.Product, error)
	ProductCheckout(payload ProductCheckoutRepositoryPayload) error
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

func (r *productRepository) GetByIds(ids []uuid.UUID) ([]entity.Product, error) {
	query := `
		SELECT id, "name", sku, category, image_url, notes, price, stock, "location", is_available, created_by, created_at, updated_at, deleted_at
			FROM public.products where id = ANY($1) AND deleted_at IS NULL
	`

	rows, err := r.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Category, &product.ImageUrl, &product.Notes, &product.Price, &product.Stock, &product.Location, &product.IsAvailable, &product.CreatedBy, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
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
		WHERE p.deleted_at IS NULL
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

	if options.SortPrice.String() != "" {
		sorting["p.price"] = options.SortPrice
	}

	if options.SortCreatedAt.String() != "" {
		sorting["p.created_at"] = options.SortCreatedAt
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

type ProductCheckoutRepositoryPayload struct {
	Customer       entity.Customer          `json:"customer"`
	ProductDetails []ProductCheckoutDetails `json:"productDetails"`
	Paid           int                      `json:"paid"`
	Change         int                      `json:"change"`
	User           entity.User              `json:"user"`
}

type ProductCheckoutDetails struct {
	ProductID string `json:"productId" validate:"uuid4,required"`
	Quantity  int    `json:"quantity" validate:"min=1,required, number"`
}

func (r *productRepository) ProductCheckout(payload ProductCheckoutRepositoryPayload) error {
	checkoutItem := make([]entity.ProductCheckoutItem, len(payload.ProductDetails))
	var cleanProducts []entity.Product
	// TODO: implement product checking
	var productIds []uuid.UUID

	for _, productDetail := range payload.ProductDetails {
		productIds = append(productIds, uuid.MustParse(productDetail.ProductID))
	}

	products, err := r.GetByIds(productIds)
	if err != nil {
		return err
	}

	// TODO: implement product checking
	// TODO: Must improve this algorithm

	if len(productIds) != len(products) {
		return commons.CustomError{Message: "product not found", Code: 404}
	}
	totalCost := 0
	for index, payloaditem := range payload.ProductDetails {
		itemCost := 0
		for _, product := range products {
			if product.ID != payloaditem.ProductID {
				continue
			} else {
				if payloaditem.Quantity > product.Stock {
					return commons.CustomError{Message: "product quantity not enough", Code: 400}
				}
				if !product.IsAvailable {
					return commons.CustomError{Message: "product quantity is not available", Code: 400}

				}
				product.Stock = product.Stock - payloaditem.Quantity
				itemCost = payloaditem.Quantity * product.Price
				checkoutItem[index].Quantity = payloaditem.Quantity
				checkoutItem[index].Amount = payloaditem.Quantity * product.Price
				checkoutItem[index].ProductID = uuid.MustParse(payloaditem.ProductID)
				cleanProducts = append(cleanProducts, product)

			}
		}
		totalCost += itemCost
	}

	// TODO: implement product calculation
	customerMoney := payload.Paid * 100
	totalCost = totalCost * 100
	if customerMoney < totalCost {
		return commons.CustomError{Message: "customer money not enough", Code: 400}
	}
	customerChange := customerMoney - totalCost

	// TODO: implement product checkout inserting data
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	log.Println("masuk sini ")
	query :=
		`
		INSERT INTO checkouts
		( customer_id, paid, "change", created_by)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at;
	`
	var checkoutID uuid.UUID
	var checkoutCreatedAt time.Time

	err = tx.QueryRow(query, payload.Customer.ID, customerMoney, customerChange, payload.User.UserID).Scan(&checkoutID, &checkoutCreatedAt)
	if err != nil {
		return err
	}

	for i := range checkoutItem {
		checkoutItem[i].CheckoutID = checkoutID
	}
	log.Println("masuk sini 2")
	// TODO: Prepare the bulk insert query
	valueStrings := make([]string, 0, len(checkoutItem))
	valueArgs := make([]interface{}, 0, len(checkoutItem)*4)
	for i, item := range checkoutItem {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, item.CheckoutID, item.ProductID, item.Quantity, item.Amount)
	}
	log.Println("masuk sini 3")
	stmt := fmt.Sprintf("INSERT INTO checkout_items (checkout_id, product_id, quantity, amount) VALUES %s", strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}
	query3 := `
		UPDATE products
			SET stock= $1, updated_at=CURRENT_TIMESTAMP
			WHERE id=$2
	`
	for _, cle := range cleanProducts {
		_, err := tx.Exec(query3, cle.Stock, cle.ID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
