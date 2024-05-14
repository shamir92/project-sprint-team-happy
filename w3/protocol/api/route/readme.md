# Routes Documentation

This document outlines the routing structure of the application, categorizing routes into public and private routes.

## Table of Contents

1. [Public Routes](#public-routes)
2. [Private Routes](#private-routes)
3. [Authentication Middleware](#authentication-middleware)

## Public Routes

Public routes are accessible to all users without any authentication requirements. These routes typically include pages like the home page, login, signup, and other information pages.

### Examples:

- `GET /`: Home page
- `GET /about`: About page
- `GET /contact`: Contact page
- `POST /login`: Login endpoint
- `POST /signup`: Signup endpoint

## Private Routes

Private routes require users to be authenticated before accessing them. These routes use the authentication middleware to ensure that only authorized users can access them. Private routes typically include user dashboard, account settings, and other user-specific features.

### Examples:

- `GET /dashboard`: User dashboard
- `GET /profile`: User profile
- `POST /account/settings`: Update account settings
- `GET /orders`: User order history

## Authentication Middleware

The authentication middleware is used to protect private routes. It checks whether the user is authenticated by verifying their authentication token or session. If the user is not authenticated, they are redirected to the login page or receive an unauthorized access response.

### Middleware Function:

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var JwtSecret = []byte("your-secret-key")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
