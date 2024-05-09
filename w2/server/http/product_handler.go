package httpserver

import (
	"eniqlostore/commons"
	"eniqlostore/internal/service"
	"fmt"
	"log"
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

	payload.CreatedBy = fmt.Sprint(r.Context().Value(currentUserRequestKey))
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
	log.Println(payload.ID)
	_, err := s.productService.UpdateProduct(payload)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}

func (s *HttpServer) handleProductDelete(w http.ResponseWriter, r *http.Request) {
	userID := fmt.Sprint(r.Context().Value(currentUserRequestKey))
	fmt.Println("shamir ->", chi.URLParam(r, "shamir"))
	err := s.productService.DeleteProduct(chi.URLParam(r, "shamir"), userID)
	if err != (commons.CustomError{}) {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "success"})
}
