package repository

import (
	"database/sql"
	"eniqlostore/internal/entity"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *customerRepository {
	return &customerRepository{
		db: db,
	}
}

func (c *customerRepository) CheckExistByPhoneNumber(phoneNumber string) (isExist bool, err error) {
	query := `
		SELECT COUNT(phone_number) FROM customers WHERE phone_number = $1
	`

	var count int
	err = c.db.QueryRow(query, phoneNumber).Scan(&count)
	isExist = count > 0

	return isExist, err
}

func (c *customerRepository) Insert(cust entity.Customer) (entity.Customer, error) {
	query := `
		INSERT INTO customers(name, phone_number) 
		VALUES($1, $2)
		RETURNING id
		`

	newCust := cust
	err := c.db.QueryRow(query, cust.Name, cust.PhoneNumber).Scan(&newCust.ID)

	return newCust, err
}
