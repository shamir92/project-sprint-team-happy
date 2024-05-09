package service

import (
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"eniqlostore/internal/repository"
	"log"
	"time"
)

type ProductServiceDeps struct {
	ProductRepository repository.IProductRepository
}

type ProductService struct {
	productRepository repository.IProductRepository
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{productRepository: deps.ProductRepository}
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedBy   string `json:"-"`
}

type CreateProductResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateProductRequest struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
}

func (s *ProductService) CreateProduct(req CreateProductRequest) (CreateProductResponse, error) {
	var resp CreateProductResponse
	product, err := entity.NewProduct(req.Name, req.SKU, req.Category, req.ImageUrl, req.Notes, req.Price, req.Stock, req.Location, req.IsAvailable, req.CreatedBy)
	if err != nil {
		return resp, err
	}

	product, err = s.productRepository.Insert(product)
	if err != nil {
		return resp, err
	}

	resp.ID = product.ID
	resp.CreatedAt = product.CreatedAt

	return resp, nil
}

func (s *ProductService) UpdateProduct(req UpdateProductRequest) (entity.Product, error) {
	_, err := s.productRepository.GetById(req.ID)
	if err != nil {
		return entity.Product{}, err
	}

	product, err := entity.NewProduct(req.Name, req.SKU, req.Category, req.ImageUrl, req.Notes, req.Price, req.Stock, req.Location, req.IsAvailable, "")
	if err != nil {
		return entity.Product{}, err
	}

	product.ID = req.ID
	err = s.productRepository.Update(product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(productId string, userId string) commons.CustomError {
	log.Println("delete product 1")
	product, err := s.productRepository.GetById(productId)
	if err != nil {
		return commons.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	log.Println("delete product 2")

	if product == (entity.Product{}) {
		return commons.CustomError{
			Message: "product not found",
			Code:    404,
		}
	}
	log.Println("delete product 3")

	if product.CreatedBy != userId {
		return commons.CustomError{
			Message: "product is not yours",
			Code:    401,
		}
	}

	log.Println("delete product 4")

	err = s.productRepository.Delete(productId)
	if err != nil {
		return commons.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}

	return commons.CustomError{}
}
