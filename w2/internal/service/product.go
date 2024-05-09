package service

import (
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"eniqlostore/internal/repository"
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

func (s *ProductService) UpdateProduct(req UpdateProductRequest, userId string) (entity.Product, error) {
	productInfo, err := s.productRepository.GetById(req.ID)
	if err != nil {
		return entity.Product{}, commons.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}

	if productInfo == (entity.Product{}) {
		return entity.Product{}, commons.CustomError{
			Message: "product not found",
			Code:    404,
		}
	}

	if productInfo.CreatedBy != userId {
		return entity.Product{}, commons.CustomError{
			Message: "product is not yours",
			Code:    401,
		}
	}

	if req.Name != "" {
		productInfo.Name = req.Name
	}
	if req.SKU != "" {
		productInfo.SKU = req.SKU
	}
	if req.Category != "" {
		productInfo.Category = req.Category
	}
	if req.ImageUrl != "" {
		productInfo.ImageUrl = req.ImageUrl
	}
	if req.Notes != "" {
		productInfo.Notes = req.Notes
	}
	if req.Price != 0 {
		productInfo.Price = req.Price
	}
	if req.Stock != 0 {
		productInfo.Stock = req.Stock
	}

	err = s.productRepository.Update(productInfo)
	if err != nil {
		return entity.Product{}, commons.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}

	return productInfo, commons.CustomError{}
}

func (s *ProductService) DeleteProduct(productId string, userId string) error {
	_, err := s.productRepository.GetById(productId)
	if err != nil {
		return err
	}

	err = s.productRepository.Delete(productId)
	if err != nil {
		return err
	}

	return nil
}
