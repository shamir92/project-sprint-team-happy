package repository

import (
	"belimang/domain/entity"
	"database/sql"
)

type IOrderRepository interface {
	Insert(entity.Order) (entity.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Insert(order entity.Order) (entity.Order, error) {
	q := `
		INSERT INTO 
			orders(user_id, user_lat, user_lon, total_price, estimated_delivery_time, state)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var newOrder = order
	row := r.db.QueryRow(q, order.UserID, order.UserLat, order.UserLon, order.TotalPrice, order.EstimatedDeliveryTime, order.State.String())

	if err := row.Scan(&newOrder.ID); err != nil {
		return newOrder, err
	}

	return newOrder, nil
}
