package service

import (
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"eniqlostore/internal/repository"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	errProductNotFound = errors.New("product not found")
)

// Apakah bisa kita hilangin ini bila gak terlalu banyak?. ngurangin2 struct.
type ProductServiceDeps struct {
	ProductRepository  repository.IProductRepository
	CustomerRepository repository.ICustomerRepository
	UserRepository     repository.IUserRepository
}

type ProductService struct {
	productRepository  repository.IProductRepository
	customerRepository repository.ICustomerRepository
	userRepository     repository.IUserRepository
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		productRepository:  deps.ProductRepository,
		customerRepository: deps.CustomerRepository,
		userRepository:     deps.UserRepository,
	}
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Notes       string `json:"notes"`
	Price       *int   `json:"price"`
	Stock       *int   `json:"stock"`
	Location    string `json:"location"`
	IsAvailable *bool  `json:"isAvailable"`
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
	Price       *int   `json:"price"`
	Stock       *int   `json:"stock"`
	Location    string `json:"location"`
	IsAvailable *bool  `json:"isAvailable"`
}

func (s *ProductService) CreateProduct(req CreateProductRequest) (CreateProductResponse, error) {
	var resp CreateProductResponse

	if req.Price == nil {
		return resp, commons.CustomError{
			Message: "price cannot be empty",
			Code:    400,
		}
	}

	if req.Stock == nil {
		return resp, commons.CustomError{
			Message: "stock cannot be empty and must be a number",
			Code:    400,
		}
	}

	if req.IsAvailable == nil {
		return resp, commons.CustomError{
			Message: "isAvailable cannot be empty and must be a boolean",
			Code:    400,
		}
	}

	product, err := entity.NewProduct(req.Name, req.SKU, req.Category, req.ImageUrl, req.Notes, *req.Price, *req.Stock, req.Location, *req.IsAvailable, req.CreatedBy)
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
	if err := uuid.Validate(req.ID); err != nil {
		return entity.Product{}, commons.CustomError{
			Message: errProductNotFound.Error(),
			Code:    404,
		}
	}

	var emptyProduct entity.Product
	if req.Price == nil {
		return emptyProduct, commons.CustomError{
			Message: "price cannot be empty",
			Code:    400,
		}
	}

	if req.Stock == nil {
		return emptyProduct, commons.CustomError{
			Message: "stock cannot be empty and must be a number",
			Code:    400,
		}
	}

	if req.IsAvailable == nil {
		return emptyProduct, commons.CustomError{
			Message: "isAvailable cannot be empty and must be a boolean",
			Code:    400,
		}
	}

	if err := entity.ValidateProductName(req.Name); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateProductSKU(req.SKU); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateProductCategory(req.Category); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateProductImageUrl(req.ImageUrl); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateNotes(req.Notes); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidatePrice(*req.Price); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateStock(*req.Stock); err != nil {
		return emptyProduct, err
	}

	if err := entity.ValidateLocation(req.Location); err != nil {
		return emptyProduct, err
	}

	productInfo, err := s.productRepository.GetById(req.ID)

	if err != nil {
		return emptyProduct, err
	}

	productInfo.Name = req.Name
	productInfo.SKU = req.SKU
	productInfo.Category = req.Category
	productInfo.Notes = req.Notes
	productInfo.ImageUrl = req.ImageUrl
	productInfo.Price = *req.Price
	productInfo.Stock = *req.Stock
	productInfo.Location = req.Location
	productInfo.IsAvailable = *req.IsAvailable

	err = s.productRepository.Update(productInfo)
	if err != nil {
		return entity.Product{}, commons.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}

	return productInfo, nil
}

func (s *ProductService) DeleteProduct(productId string, userId string) error {

	if err := uuid.Validate(productId); err != nil {
		return commons.CustomError{
			Message: errProductNotFound.Error(),
			Code:    404,
		}
	}

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

type GetProductsRequest struct {
	Limit         string `json:"limit"`
	Offset        string `json:"offset"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	SKU           string `json:"sku"`
	IsAvailable   string `json:"isAvailable"`
	InStock       string `json:"inStock"`
	SortPrice     string `json:"price"`
	SortCreatedAt string `json:"createdAt"`
}

func (s *ProductService) GetProducts(req GetProductsRequest) ([]entity.Product, error) {
	var options []entity.FindProductOptionBuilder

	var limit, offset = 5, 0

	if l, err := strconv.Atoi(req.Limit); err == nil && l > 0 {
		limit = l
	}

	if o, err := strconv.Atoi(req.Offset); err == nil && o > 0 {
		offset = o
	}

	options = append(options, entity.WithOffsetAndLimit(offset, limit))

	if req.ID != "" {
		options = append(options, entity.WithProductID(req.ID))
	}

	if req.Name != "" {
		options = append(options, entity.WithProductName(req.Name))
	}

	if req.Category != "" {
		options = append(options, entity.WithProductCategory(req.Category))
	}

	if req.SKU != "" {
		options = append(options, entity.WithProductSKU(req.SKU))
	}

	if isAvailable, err := strconv.ParseBool(req.IsAvailable); err == nil {
		options = append(options, entity.WithIsAvailable(&isAvailable))
	}

	if inStock, err := strconv.ParseBool(req.InStock); err == nil {
		options = append(options, entity.WithInStock(&inStock))
	}

	if req.SortPrice == entity.DESC.String() {
		options = append(options, entity.WithSortPrice(entity.DESC))
	} else if req.SortPrice == entity.ASC.String() {
		options = append(options, entity.WithSortPrice(entity.ASC))
	}

	if req.SortCreatedAt == entity.DESC.String() {
		options = append(options, entity.WithSortCreatedAt(entity.DESC))
	} else if req.SortCreatedAt == entity.ASC.String() {
		options = append(options, entity.WithSortCreatedAt(entity.ASC))
	} else {
		options = append(options, entity.WithSortCreatedAt(entity.DESC))
	}

	return s.productRepository.Find(options...)
}

type ProductCheckoutRequest struct {
	CustomerID     string                              `json:"customerId" validate:"uuid4,required, number"`
	ProductDetails []repository.ProductCheckoutDetails `json:"productDetails"`
	Paid           int                                 `json:"paid" validate:"min=1,required, number"`
	Change         *int                                `json:"change" validate:"min=0,required, number"`
	UserID         string                              `json:"userId" validate:"uuid4,required"`
}

func (s *ProductService) ProductCheckout(payload ProductCheckoutRequest) error {
	if err := uuid.Validate(payload.CustomerID); err != nil {
		return commons.CustomError{
			Message: "customer's id is invalid",
			Code:    400,
		}
	}

	if payload.Change == nil {
		return commons.CustomError{
			Message: "change cannot be empty",
			Code:    400,
		}
	} else if *payload.Change < 0 {
		return commons.CustomError{
			Message: "change should be a number with minum 0 as value",
			Code:    400,
		}
	}

	if len(payload.ProductDetails) == 0 {
		return commons.CustomError{
			Message: "productDetails must be an array and cannot be empty",
			Code:    400,
		}
	}

	for _, item := range payload.ProductDetails {
		if item.ProductID == "" {
			return commons.CustomError{
				Message: "product's id is invalid",
				Code:    400,
			}
		}

		if err := uuid.Validate(item.ProductID); err != nil {
			return commons.CustomError{
				Message: "product not found",
				Code:    404,
			}
		}

		if item.Quantity < 1 {
			return commons.CustomError{
				Message: "item's quantity must be at least 1",
				Code:    400,
			}
		}
	}

	cust, err := s.customerRepository.GetById(payload.CustomerID)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetById(payload.UserID)
	if err != nil {
		return err
	}

	err = s.productRepository.ProductCheckout(repository.ProductCheckoutRepositoryPayload{
		Customer:       cust,
		ProductDetails: payload.ProductDetails,
		Paid:           payload.Paid,
		Change:         *payload.Change,
		User:           user,
	})

	if err != nil {
		return err
	}

	return nil
}
