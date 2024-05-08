package httpserver

import (
	"eniqlostore/internal/service"
	"net/http"
)

func (s *HttpServer) handleStaffCreate(w http.ResponseWriter, r *http.Request) {
	payload := service.CreateStaffRequest{}

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	newStaff, err := s.userService.UserCreate(payload)

	if err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"data": newStaff})
}

func (s *HttpServer) handleStaffLogin(w http.ResponseWriter, r *http.Request) {
	payload := service.UserLoginRequest{}

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	stafLogedIn, err := s.userService.UserLogin(payload)

	if err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"data": stafLogedIn})
}
