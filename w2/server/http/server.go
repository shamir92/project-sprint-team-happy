package httpserver

import (
	"database/sql"
	"encoding/json"
	"eniqlostore/internal/auth"
	"eniqlostore/internal/repository"
	"eniqlostore/internal/service"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ServerOpts struct {
	Addr string
	DB   *sql.DB
}

type HttpServer struct {
	addr            string // todo: change to HttpServerConfig
	userService     *service.UserService
	productService  *service.ProductService
	customerService *service.CustomerService
	checkoutService *service.ProductCheckoutService
	tokenManager    auth.AuthJwtTokenManager
}

func New(opts ServerOpts) *HttpServer {
	jwtTokenManager := auth.NewJwt()
	userRepo := repository.NewUserRepository(opts.DB)
	userService := service.NewUserService(service.UserServiceDeps{
		UserRepository:   userRepo,
		AuthTokenManager: jwtTokenManager,
		PasswordHash:     auth.NewBcryptPasswordHash(),
	})

	custRepo := repository.NewCustomerRepository(opts.DB)
	custService := service.NewCustomerService(custRepo)

	productRepo := repository.NewProductRepository(opts.DB)
	productService := service.NewProductService(service.ProductServiceDeps{
		ProductRepository:  productRepo,
		CustomerRepository: custRepo,
		UserRepository:     userRepo,
	})

	checkoutRepo := repository.NewProductCheckoutRepository(opts.DB)
	checkoutService := service.NewProductCheckoutService(checkoutRepo)

	return &HttpServer{
		addr:            opts.Addr,
		userService:     userService,
		productService:  productService,
		customerService: custService,
		tokenManager:    jwtTokenManager,
		checkoutService: checkoutService,
	}
}

func (s *HttpServer) Server() *http.Server {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	router.Use(middleware.Heartbeat("/ping"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey, What's Up!"))
	})

	// private routes
	router.Route("/v1", func(v1 chi.Router) {
		// public routes
		v1.Route("/staff", func(r chi.Router) {
			r.Post("/register", s.handleStaffCreate)
			r.Post("/login", s.handleStaffLogin)
		})

		v1.Get("/product/customer", s.handleSearchProducts)

		v1.Post("/ping", s.handlePing)

		v1.Route("/product", func(productRouter chi.Router) {
			productRouter.Use(s.AuthMiddleware)
			productRouter.Get("/", s.handleProductBrowse)
			productRouter.Post("/", s.handleProductCreate)
			productRouter.Put("/{productId}", s.handleProductEdit)
			productRouter.Delete("/{productId}", s.handleProductDelete)
			productRouter.Post("/checkout", s.handleProductCheckout)
			productRouter.Get("/checkout/history", s.handleGetProductCheckoutHistories)
		})

		v1.Route("/customer", func(custRouter chi.Router) {
			custRouter.Use(s.AuthMiddleware)
			custRouter.Post("/register", s.handleCreateCustomer)
			custRouter.Get("/", s.handleGetCustomers)
		})

	})

	return &http.Server{
		Addr:    s.addr,
		Handler: router,
	}
}

func (s *HttpServer) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) error {
	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)

	return nil
}

func (s *HttpServer) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
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
