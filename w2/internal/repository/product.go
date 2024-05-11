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
	GetByIds(ids []string) ([]entity.Product, error)
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

func (r *productRepository) GetByIds(ids []string) ([]entity.Product, error) {
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
	var productIds []string

	for _, productDetail := range payload.ProductDetails {
		id := uuid.MustParse(productDetail.ProductID)
		productIds = append(productIds, id.String())
	}

	products, err := r.GetByIds(productIds)
	if err != nil {
		return err
	}

	if len(productIds) != len(products) {
		return commons.CustomError{Message: "product not found", Code: 404}
	}

	productsById := make(map[string]entity.Product, len(products))

	for _, p := range products {
		productsById[p.ID] = p
	}

	totalCost := 0
	for index, payloaditem := range payload.ProductDetails {
		itemCost := 0

		product, ok := productsById[payloaditem.ProductID]
		if !ok {
			return commons.CustomError{
				Message: "product not found",
				Code:    404,
			}
		}

		if payloaditem.Quantity > product.Stock {
			return commons.CustomError{Message: "product quantity not enough", Code: 400}
		}

		if !product.IsAvailable {
			return commons.CustomError{Message: "product is not available yet", Code: 400}
		}

		product.Stock = product.Stock - payloaditem.Quantity
		itemCost = payloaditem.Quantity * product.Price
		checkoutItem[index].Quantity = payloaditem.Quantity
		checkoutItem[index].Amount = payloaditem.Quantity * product.Price
		checkoutItem[index].ProductID = uuid.MustParse(payloaditem.ProductID)

		totalCost += itemCost
	}

	customerMoney := payload.Paid
	if customerMoney < totalCost {
		return commons.CustomError{Message: "customer money not enough", Code: 400}
	}

	customerChange := customerMoney - totalCost
	if customerChange != payload.Change {
		return commons.CustomError{Message: "change is not valid", Code: 400}
	}

	err = r.transaction(func(tx *sql.Tx) error {
		query :=
			`
		INSERT INTO checkouts
		(customer_id, paid, "change", created_by)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at;
	`
		var checkoutID uuid.UUID
		var checkoutCreatedAt time.Time

		err = tx.QueryRow(query, payload.Customer.ID, customerMoney, customerChange, payload.User.UserID).Scan(&checkoutID, &checkoutCreatedAt)
		if err != nil {
			log.Printf("transaction insert checkouts: %v", err)
			return err
		}

		for i := range checkoutItem {
			checkoutItem[i].CheckoutID = checkoutID
		}

		// TODO: Prepare the bulk insert query
		var insertCheckoutItemValues strings.Builder
		var queryValuesForUpdateStock strings.Builder
		insertCheckoutItemArguments := make([]interface{}, 0, len(checkoutItem)*4)
		updateProductsStockArgs := make([]interface{}, 0, len(checkoutItem)*4)

		for i, item := range checkoutItem {
			if i < len(checkoutItem)-1 {
				insertCheckoutItemValues.WriteString(
					fmt.Sprintf("($%d, $%d, $%d, $%d),", i*4+1, i*4+2, i*4+3, i*4+4),
				)

				queryValuesForUpdateStock.WriteString(
					fmt.Sprintf("($%d, $%d), ", i*2+1, i*2+2),
				)
			} else {
				// Last element doesn't need comma at the end
				insertCheckoutItemValues.WriteString(
					fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4),
				)

				queryValuesForUpdateStock.WriteString(
					fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2),
				)

			}
			insertCheckoutItemArguments = append(insertCheckoutItemArguments, item.CheckoutID, item.ProductID, item.Quantity, item.Amount)
			updateProductsStockArgs = append(updateProductsStockArgs, item.ProductID, item.Quantity)
		}

		insertCheckoutItemsStmt := fmt.Sprintf("INSERT INTO checkout_items (checkout_id, product_id, quantity, amount) VALUES %s", insertCheckoutItemValues.String())
		_, err = tx.Exec(insertCheckoutItemsStmt, insertCheckoutItemArguments...)

		if err != nil {
			log.Printf("transaction insert checkout_items: %v", err)
			return err
		}

		updateProductsStokStmt := fmt.Sprintf(`
		UPDATE products p
		SET
			stock = p.stock - checkout_items.quantity::smallint
		FROM (VALUES %s) AS checkout_items(product_id, quantity)
		WHERE p.id = checkout_items.product_id::uuid
		`, queryValuesForUpdateStock.String())

		_, err = tx.Exec(updateProductsStokStmt, updateProductsStockArgs...)

		if err != nil {
			log.Printf("transaction update products: %v", err)
			log.Printf("transaction update products query: %s", updateProductsStokStmt)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type fn func(tx *sql.Tx) error

func (r *productRepository) transaction(cb fn) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	err = cb(tx)

	if err != nil {
		if rbbErr := tx.Rollback(); rbbErr != nil {
			return fmt.Errorf("tx transaction err: %v, rollback err: %v", err, rbbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
