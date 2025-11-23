// Package server: defines the api server functionality.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
)

type apiServer struct {
	prodSvc adapters.ProductService
	logger  *slog.Logger
}

type APIServerOptions func(*apiServer)

func WithLogger(logger *slog.Logger) APIServerOptions {
	return func(a *apiServer) {
		a.logger = logger
	}
}

func NewAPIServer(prodSvc adapters.ProductService, opts ...APIServerOptions) adapters.APIServer {
	apiServer := &apiServer{
		prodSvc: prodSvc,
	}

	for _, opt := range opts {
		opt(apiServer)
	}

	if apiServer.logger == nil {
		apiServer.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return apiServer
}

func (a *apiServer) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.HealthCheck)
	mux.HandleFunc("/product", a.ListProducts)

	return a.configureCorsMiddleware(mux)
}

func (a *apiServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("in good health"))
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *apiServer) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := a.prodSvc.ListProducts()
	if err != nil {
		a.logger.Error(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	jsonProds, err := json.Marshal(products)
	if err != nil {
		a.logger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonProds)
	if err != nil {
		a.logger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *apiServer) configureCorsMiddleware(h http.Handler) http.Handler {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		h.ServeHTTP(w, r)
	})

	return hf
}
