package main

import "net/http"

func (s *server) handleCreateStaff(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating new staff"))
}
