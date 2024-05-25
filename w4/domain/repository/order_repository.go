package repository

import (
	"belimang/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type IOrderRepository interface {
	InsertEstimateOrder(in InsertEstimateOrderPayload) (entity.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) insertOrder(order entity.Order, tx *sql.Tx) (entity.Order, error) {
	q := `
		INSERT INTO 
			orders(user_id, user_lat, user_lon, total_price, estimated_delivery_time, state)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var newOrder = order
	row := tx.QueryRow(q, order.UserID, order.UserLat, order.UserLon, order.TotalPrice, order.EstimatedDeliveryTime, order.State.String())

	if err := row.Scan(&newOrder.ID); err != nil {
		return newOrder, err
	}

	return newOrder, nil
}

func (r *orderRepository) insertOrderItems(orderID uuid.UUID, items []entity.OrderItem, tx *sql.Tx) error {
	q := `
		INSERT INTO order_items(order_id, order_item_id, price, quantity, amount)
	`

	var values = make([]any, 0)
	var valuesParams = make([]string, 0)
	for i, item := range items {
		var totalParams = 5
		var value = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*totalParams+1, i*totalParams+2, i*totalParams+3, i*totalParams+4, i*totalParams+5)
		valuesParams = append(valuesParams, value)
		values = append(values, orderID, item.ItemID, item.Price, item.Quantity, item.Amount)
	}

	q += fmt.Sprintf("\nVALUES %s", strings.Join(valuesParams, ","))

	_, err := tx.Exec(q, values...)

	if err != nil {
		return err
	}

	return nil
}

type InsertEstimateOrderPayload struct {
	Order      entity.Order
	OrderItems []entity.OrderItem
}

func (r *orderRepository) InsertEstimateOrder(in InsertEstimateOrderPayload) (entity.Order, error) {
	var newOrder entity.Order
	err := r.startTrx(func(tx *sql.Tx) error {
		createdOrder, err := r.insertOrder(in.Order, tx)

		if err != nil {
			return err
		}
		newOrder = createdOrder

		err = r.insertOrderItems(newOrder.ID, in.OrderItems, tx)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return newOrder, err
	}

	return newOrder, nil
}

func (r *orderRepository) startTrx(f func(tx *sql.Tx) error) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		log.Printf("ERROR | startTrx | %v\n", err)
		return err
	}

	err = f(tx)

	if err != nil {
		log.Printf("ERROR | startTrx | %v\n", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("ERROR | startTrx | rollback a transaction failed:  %v\n", rbErr)
			return fmt.Errorf("ERROR | startTrx | rollback error: %v", rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("ERROR | startTrx | commit a transaction failed:  %v\n", err)
		return err
	}

	return nil
}
