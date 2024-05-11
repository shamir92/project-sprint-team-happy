package repository

import (
	"database/sql"
	"eniqlostore/internal/entity"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type ProductCheckoutIncludeItem struct {
	entity.ProductCheckout
	Item entity.ProductCheckoutItem
}

type IProductCheckoutRepository interface {
	Find(opts ...entity.FindCheckoutHistoryBuilder) ([]entity.ProductCheckout, error)
	FindItemByCheckoutIds(checkOutIds []string) ([]entity.ProductCheckoutItem, error)
}

type checkoutRepository struct {
	db *sql.DB
}

func NewProductCheckoutRepository(db *sql.DB) *checkoutRepository {
	return &checkoutRepository{
		db: db,
	}
}

func (cp *checkoutRepository) Find(opts ...entity.FindCheckoutHistoryBuilder) ([]entity.ProductCheckout, error) {
	options := &entity.FindCheckoutHistory{}
	for _, o := range opts {
		o(options)
	}

	var whereQueries []string
	values := []interface{}{
		options.Limit,
		options.Offset,
	}

	// Filter
	if options.CustomerID != "" {
		values = append(values, options.CustomerID)
		whereQueries = append(whereQueries, fmt.Sprintf("c.customer_id = $%d", len(values)))
	}

	// Sorting
	sort := make(map[string]entity.SortType)
	if options.SortCreatedAt.String() != "" {
		sort["c.created_at"] = options.SortCreatedAt
	}

	whereStmt := "\n"
	if len(whereQueries) > 0 {
		whereStmt = fmt.Sprintf("WHERE %s", strings.Join(whereQueries, " AND "))
	}

	var sortStmtBuilder strings.Builder
	countSort := 0
	totalSort := len(sort)
	for key, val := range sort {
		countSort++

		if countSort == totalSort {
			sortStmtBuilder.WriteString(
				fmt.Sprintf("%s %s", key, val),
			)
		} else {
			sortStmtBuilder.WriteString(
				fmt.Sprintf("%s %s,", key, val),
			)
		}
	}

	var sortQuery = "\n"
	if sortStmtBuilder.String() != "" {
		sortQuery = fmt.Sprintf("ORDER BY %s", sortStmtBuilder.String())
	}

	finalQuery := fmt.Sprintf(`
		SELECT 
			c.id, c.customer_id, c.paid, c.change, c.created_at
		FROM checkouts c
		%s
		%s
		LIMIT $1 OFFSET $2
	`, whereStmt, sortQuery)

	rows, err := cp.db.Query(finalQuery, values...)

	var checkoutHistories []entity.ProductCheckout
	if err != nil {
		log.Printf("repo product checkout Find query %v", err)
		return checkoutHistories, err
	}

	defer rows.Close()

	for rows.Next() {
		var pc entity.ProductCheckout
		err := rows.Scan(
			&pc.CheckoutID,
			&pc.CustomerID,
			&pc.Paid,
			&pc.Change,
			&pc.CreatedAt,
		)

		if err != nil {
			log.Printf("repo product checkout Find scan: %v", err)
			return checkoutHistories, err
		}

		checkoutHistories = append(checkoutHistories, pc)
	}

	return checkoutHistories, nil
}

func (cp *checkoutRepository) FindItemByCheckoutIds(checkoutIds []string) ([]entity.ProductCheckoutItem, error) {
	rows, err := cp.db.Query(`
		SELECT ci.checkout_id, ci.product_id, ci.quantity
		FROM checkout_items ci
		WHERE ci.checkout_id = ANY($1)
	`, pq.Array(checkoutIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []entity.ProductCheckoutItem
	for rows.Next() {
		var item entity.ProductCheckoutItem

		if err := rows.Scan(&item.CheckoutID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
