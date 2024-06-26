package httpserver

import (
	"eniqlostore/internal/service"
	"errors"
	"net/http"
)

func (s *HttpServer) handleStaffCreate(w http.ResponseWriter, r *http.Request) {
	payload := service.CreateStaffRequest{}

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	newStaff, err := s.userService.UserCreate(payload)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusCreated, map[string]any{"data": newStaff})
}

func (s *HttpServer) handleStaffLogin(w http.ResponseWriter, r *http.Request) {
	payload := service.UserLoginRequest{}

	if err := s.decodeJSON(w, r, &payload); err != nil {
		s.errorBadRequest(w, r, err)
		return
	}

	if payload.Password == nil || payload.PhoneNumber == nil {
		s.errorBadRequest(w, r, errors.New("invalid payload"))
		return
	} else if *payload.Password == "" || *payload.PhoneNumber == "" {
		s.errorBadRequest(w, r, errors.New("invalid payload"))
		return
	}

	stafLogedIn, err := s.userService.UserLogin(payload)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.writeJSON(w, r, http.StatusOK, map[string]any{"data": stafLogedIn})
}
