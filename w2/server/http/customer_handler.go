package httpserver

import (
	"eniqlostore/internal/service"
	"net/http"
)

func (s *HttpServer) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	body := service.CreateCustomerRequest{}

	if err := s.decodeJSON(w, r, &body); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	cust, err := s.customerService.CreateCustomer(body)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"data": cust})
}

func (s *HttpServer) handleGetCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	customers, err := s.customerService.GetCustomers(service.GetCustomerRequst{
		Name:        query.Get("name"),
		PhoneNumber: query.Get("phoneNumber"),
	})

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"data": customers})
}
