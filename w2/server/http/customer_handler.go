package httpserver

import (
	"eniqlostore/internal/service"
	"net/http"
)

func (s *HttpServer) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	body := service.CreateCustomerRequest{}

	if err := s.decodeJSON(w, r, &body); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	cust, err := s.customerService.CreateCustomer(body)

	if err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"data": cust})
}