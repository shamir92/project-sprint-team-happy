package repository

import (
	"belimang/domain/entity"
	"belimang/internal/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type IOrderRepository interface {
	InsertEstimateOrder(in InsertEstimateOrderPayload) (entity.Order, error)
	FindOrderByID(orderID string) (entity.Order, error)
	Update(entity.Order) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindOrderByID(orderID string) (entity.Order, error) {
	q := `
		SELECT 
			id, user_id, user_lon, user_lat, total_price, estimated_delivery_time
		FROM orders
		WHERE id = $1
		`

	var order entity.Order

	err := r.db.QueryRow(q, orderID).Scan(&order.ID, &order.UserID, &order.UserLon, &order.UserLat, &order.TotalPrice, &order.EstimatedDeliveryTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, helper.CustomError{
				Code:    404,
				Message: "order not found",
			}
		}

		log.Printf("ERROR | FindOrderByID() | %v\n", err)
		return order, err
	}

	return order, nil
}

func (r *orderRepository) Update(order entity.Order) error {
	updateQuery := `UPDATE orders SET state = $1 WHERE id = $2`

	result, err := r.db.Exec(updateQuery, order.State, order.ID)

	if err != nil {
		log.Printf("ERROR | OrderRepository.Update() | %v\n", err)
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		log.Printf("ERROR | OrderRepository.Update() | %v\n", err)
		return err
	}

	if affected != 1 {
		errMsg := fmt.Sprintf("ERROR | OrderRepository.Update() | failed to update rows, expected 1 row to be affected, but %d rows were affected", affected)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	return nil
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
