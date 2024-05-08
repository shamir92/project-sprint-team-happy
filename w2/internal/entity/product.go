package entity

import (
	"eniqlostore/commons"
	"fmt"
	"time"
)

type Product struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	SKU         string     `json:"sku"`
	Category    string     `json:"category"`
	ImageUrl    string     `json:"imageUrl"`
	Notes       string     `json:"notes"`
	Price       int        `json:"price"`
	Stock       int        `json:"stock"`
	Location    string     `json:"location"`
	IsAvailable bool       `json:"isAvailable"`
	CreatedBy   string     `json:"-"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

func validateStock(stock int) error {
	const MIN_STOCK = 0
	const MAX_STOCK = 10000
	if stock < MIN_STOCK {
		return commons.CustomError{Message: fmt.Sprintf("stock must be greater or equal than %d", MIN_STOCK), Code: 400}
	}

	if stock > MAX_STOCK {
		return commons.CustomError{Message: fmt.Sprintf("stock must be less or equal than %d", MAX_STOCK), Code: 400}
	}

	return nil
}

func validatePrice(price int) error {
	const MIN_PRICE = 1

	if price < MIN_PRICE {
		return commons.CustomError{Message: fmt.Sprintf("price must be greater or equal than %d", MIN_PRICE), Code: 400}
	}

	return nil
}

func validateNotes(notes string) error {
	const MIN_NOTES = 1
	const MAX_NOTES = 200

	if len(notes) < MIN_NOTES || len(notes) > MAX_NOTES {
		return commons.CustomError{
			Message: fmt.Sprintf("notes should be between %d and %d characters long", MIN_NOTES, MAX_NOTES), Code: 400}
	}

	return nil
}

func validateLocation(location string) error {
	const MIN_LOCATION = 1
	const MAX_LOCATION = 200

	if len(location) < MIN_LOCATION || len(location) > MAX_LOCATION {
		return commons.CustomError{
			Message: fmt.Sprintf("location should be between %d and %d characters long", MIN_LOCATION, MAX_LOCATION), Code: 400}
	}

	return nil
}

func NewProduct(
	name string,
	sku string,
	category string,
	imageUrl string,
	notes string,
	price int,
	stock int,
	location string,
	isAvailable bool,
	createdBy string,
) (Product, error) {
	var product Product

	if err := validateStock(stock); err != nil {
		return product, err
	}

	if err := validatePrice(price); err != nil {
		return product, err
	}

	if err := validateNotes(notes); err != nil {
		return product, err
	}

	if err := validateLocation(location); err != nil {
		return product, err
	}

	product.Name = name
	product.SKU = sku
	product.Category = category
	product.ImageUrl = imageUrl
	product.Notes = notes
	product.Price = price
	product.Stock = stock
	product.Location = location
	product.IsAvailable = isAvailable
	product.CreatedBy = createdBy

	return product, nil
}
