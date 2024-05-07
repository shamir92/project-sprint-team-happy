package main

import (
	"encoding/json"
	"eniqlostore/internal/service"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type server struct {
	service *service.Service
	router  *chi.Mux
}

func newServer(service *service.Service) *server {
	s := &server{
		service: service,
	}

	s.router = chi.NewRouter()

	s.routes()

	return s
}

func (s *server) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) error {
	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)

	return nil
}

func (s *server) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)

	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		/*
			A json.InvalidUnmarshalError error will be returned if we pass something
			that is not a non-nil pointer to Decode(). We catch this and panic,
			rather than returning an error to our handler.
		*/
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	return nil
}
