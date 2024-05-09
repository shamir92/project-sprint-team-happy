package httpserver

import (
	"context"
	"eniqlostore/commons"
	"net/http"
	"strings"
)

type CurrentUserRequest string

const (
	currentUserRequestKey CurrentUserRequest = "user"
)

func (s *HttpServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")

		var errUnauthorized commons.CustomError = commons.CustomError{
			Message: "Unauthorized",
			Code:    http.StatusBadRequest,
		}

		if authHeader == "" {
			s.errorResponse(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		authHeaderValues := strings.Split(authHeader, "Bearer ")
		if len(authHeaderValues) < 2 {
			s.errorResponse(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		token := authHeaderValues[1]

		if claim, err := s.tokenManager.GetClaim(token); err == nil {
			ctx := context.WithValue(r.Context(), currentUserRequestKey, claim.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			errUnauthorized.Message = err.Error()
			s.errorResponse(w, r, http.StatusUnauthorized, errUnauthorized)
		}

	})
}
