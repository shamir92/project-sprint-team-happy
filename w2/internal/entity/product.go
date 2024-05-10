package entity

import (
	"eniqlostore/commons"
	"fmt"
	"regexp"
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

func ValidateProductName(name string) error {
	const MIN_LENGTH = 1
	const MAX_LENGTH = 30

	if length := len(name); length < MIN_LENGTH || length > MAX_LENGTH {
		return commons.CustomError{Message: fmt.Sprintf("name cannot be empty and must be between %d and %d characters long", MIN_LENGTH, MAX_LENGTH), Code: 400}
	}

	return nil
}

func ValidateProductSKU(sku string) error {
	const MIN_LENGTH = 1
	const MAX_LENGTH = 30

	if length := len(sku); length < MIN_LENGTH || length > MAX_LENGTH {
		return commons.CustomError{Message: fmt.Sprintf("sku cannot be empty and must be between %d and %d characters long", MIN_LENGTH, MAX_LENGTH), Code: 400}
	}

	return nil
}

func ValidateStock(stock int) error {
	const MIN_STOCK = 1
	const MAX_STOCK = 10000
	if stock < MIN_STOCK {
		return commons.CustomError{Message: fmt.Sprintf("stock must be greater or equal than %d", MIN_STOCK), Code: 400}
	}

	if stock > MAX_STOCK {
		return commons.CustomError{Message: fmt.Sprintf("stock must be less or equal than %d", MAX_STOCK), Code: 400}
	}

	return nil
}

func ValidatePrice(price int) error {
	const MIN_PRICE = 1

	if price < MIN_PRICE {
		return commons.CustomError{Message: fmt.Sprintf("price must be greater or equal than %d", MIN_PRICE), Code: 400}
	}

	return nil
}

func ValidateNotes(notes string) error {
	const MIN_NOTES = 1
	const MAX_NOTES = 200

	if len(notes) < MIN_NOTES || len(notes) > MAX_NOTES {
		return commons.CustomError{
			Message: fmt.Sprintf("notes should be between %d and %d characters long", MIN_NOTES, MAX_NOTES), Code: 400}
	}

	return nil
}

func ValidateProductCategory(category string) error {
	if category == "" {
		return commons.CustomError{
			Message: "category cannot be empty",
			Code:    400,
		}
	}

	return nil
}

func ValidateProductImageUrl(rawURL string) error {
	pattern := `^(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the URL matches the pattern
	if !regex.MatchString(rawURL) {
		return commons.CustomError{
			Message: "imageUrl should be valid url",
			Code:    400,
		}
	}

	return nil
}

func ValidateLocation(location string) error {
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

	if err := ValidateProductName(name); err != nil {
		return product, err
	}

	if err := ValidateProductSKU(sku); err != nil {
		return product, err
	}

	if err := ValidateProductCategory(category); err != nil {
		return product, err
	}

	if err := ValidateProductImageUrl(imageUrl); err != nil {
		return product, err
	}

	if err := ValidateStock(stock); err != nil {
		return product, err
	}

	if err := ValidatePrice(price); err != nil {
		return product, err
	}

	if err := ValidateNotes(notes); err != nil {
		return product, err
	}

	if err := ValidateLocation(location); err != nil {
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
