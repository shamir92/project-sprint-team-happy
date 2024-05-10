package httpserver

import (
	"eniqlostore/internal/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *HttpServer) handleProductBrowse(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	products, err := s.productService.GetProducts(service.GetProductsRequest{
		Limit:         query.Get("limit"),
		Offset:        query.Get("offset"),
		ID:            query.Get("id"),
		Name:          query.Get("name"),
		IsAvailable:   query.Get("isAvailable"),
		Category:      query.Get("category"),
		SKU:           query.Get("sku"),
		SortPrice:     query.Get("price"),
		InStock:       query.Get("inStock"),
		SortCreatedAt: query.Get("createdAt"),
	})

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success", "data": products})
}

func (s *HttpServer) handleProductCreate(w http.ResponseWriter, r *http.Request) {

	var payload service.CreateProductRequest

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	payload.CreatedBy = fmt.Sprint(r.Context().Value(currentUserRequestKey))
	product, err := s.productService.CreateProduct(payload)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"message": "success", "data": product})
}

func (s *HttpServer) handleProductEdit(w http.ResponseWriter, r *http.Request) {
	userID := fmt.Sprint(r.Context().Value(currentUserRequestKey))

	var payload service.UpdateProductRequest

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	payload.ID = chi.URLParam(r, "productId")
	_, err := s.productService.UpdateProduct(payload, userID)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}

func (s *HttpServer) handleProductDelete(w http.ResponseWriter, r *http.Request) {
	userID := fmt.Sprint(r.Context().Value(currentUserRequestKey))
	err := s.productService.DeleteProduct(chi.URLParam(r, "productId"), userID)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}

func (s *HttpServer) handleSearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	products, err := s.productService.GetProducts(service.GetProductsRequest{
		Limit:         query.Get("limit"),
		Offset:        query.Get("offset"),
		Name:          query.Get("name"),
		IsAvailable:   "true",
		Category:      query.Get("category"),
		SKU:           query.Get("sku"),
		InStock:       query.Get("inStock"),
		SortPrice:     query.Get("price"),
		SortCreatedAt: query.Get("createdAt"),
	})

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success", "data": products})
}
