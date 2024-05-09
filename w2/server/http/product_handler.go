package httpserver

import (
	"eniqlostore/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *HttpServer) handleProductBrowse(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}

func (s *HttpServer) handleProductCreate(w http.ResponseWriter, r *http.Request) {
	var payload service.CreateProductRequest

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	payload.CreatedBy = "44d300ce-c62c-421b-a432-78c825da877a"
	product, err := s.productService.CreateProduct(payload)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"message": "success", "data": product})
}

func (s *HttpServer) handleProductEdit(w http.ResponseWriter, r *http.Request) {
	var payload service.UpdateProductRequest

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	payload.ID = chi.URLParam(r, "productId")
	_, err := s.productService.UpdateProduct(payload)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}

func (s *HttpServer) handleProductDelete(w http.ResponseWriter, r *http.Request) {
	err := s.productService.DeleteProduct(chi.URLParam(r, "productId"))
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}
