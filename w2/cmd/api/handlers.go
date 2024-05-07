package main

import (
	"eniqlostore/internal/service"
	"net/http"
)

func (s *server) handleStaffCreate(w http.ResponseWriter, r *http.Request) {
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
