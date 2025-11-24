// Package server: defines the api server functionality.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/sunimalherath/orderfoodonline/internal/app/utils"
	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type apiServer struct {
	prodSvc  adapters.ProductService
	orderSvc adapters.OrderService
	apiKey   string
	logger   *slog.Logger
}

type APIServerOptions func(*apiServer)

func WithLogger(logger *slog.Logger) APIServerOptions {
	return func(a *apiServer) {
		a.logger = logger
	}
}

func NewAPIServer(prodSvc adapters.ProductService, orderSvc adapters.OrderService, opts ...APIServerOptions) adapters.APIServer {
	apiServer := &apiServer{
		prodSvc:  prodSvc,
		orderSvc: orderSvc,
		apiKey:   utils.GetEnvVar(constants.APIKey, "apitest"),
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

	mux.HandleFunc("GET /health", a.HealthCheck)
	mux.HandleFunc("GET /product", a.ListProducts)
	mux.HandleFunc("GET /product/{productId}", a.FindProductByID)
	mux.HandleFunc("POST /order", a.PlaceAnOrder)

	return a.configureCorsMiddleware(a.authAPIkeyMiddleware(mux))
}

func (a *apiServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	a.logger.Info(constants.CheckHealth)
	a.writeJSONResponse(w, http.StatusOK, constants.SUCCESS, constants.GoodHealth, nil)
}

func (a *apiServer) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ActiveDuration)
	defer cancel()

	products, err := a.prodSvc.ListProducts(ctx)
	if err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, http.StatusInternalServerError, constants.FAILURE, constants.RetrievalFailed, nil)

		return

	}

	a.writeJSONResponse(w, http.StatusOK, constants.SUCCESS, constants.ProductsRcvd, products)
}

func (a *apiServer) FindProductByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ActiveDuration)
	defer cancel()

	prodID, err := strconv.Atoi(r.PathValue("productId"))
	if err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, http.StatusBadRequest, constants.FAILURE, constants.InvalidProdID, nil)

		return
	}

	product, err := a.prodSvc.FindProductByID(ctx, int64(prodID))
	if err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, http.StatusNotFound, constants.FAILURE, constants.ProdNotFound, nil)

		return
	}

	if product == nil {
		a.logger.Warn("product not found for product Id", slog.Int("productId", prodID))

		msg := fmt.Sprintf("product not for product ID: %d", prodID)
		a.writeJSONResponse(w, http.StatusInternalServerError, constants.FAILURE, msg, nil)

		return
	}

	msg := fmt.Sprintf("product retrieved for productId: %d", prodID)

	a.writeJSONResponse(w, http.StatusOK, constants.SUCCESS, msg, product)
}

func (a *apiServer) PlaceAnOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ActiveDuration)
	defer cancel()

	var orderReq entities.OrderReq

	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, http.StatusBadRequest, constants.FAILURE, constants.InvalidRequest, nil)

		return
	}

	if err := orderReq.Validate(); err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, http.StatusUnprocessableEntity, constants.FAILURE, constants.ValidationFailed, nil)

		return
	}

	order, err := a.orderSvc.PlaceAnOrder(ctx, orderReq)
	if err != nil {
		a.logger.Error(err.Error())

		a.writeJSONResponse(w, getStatusCode(err), constants.FAILURE, err.Error(), nil)

		return
	}

	a.writeJSONResponse(w, http.StatusOK, constants.SUCCESS, constants.OrderPlaced, order)
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

func (a *apiServer) authAPIkeyMiddleware(h http.Handler) http.Handler {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("api_key")

		if key == "" {
			a.logger.Error(constants.MissingAPIkey)
			a.writeJSONResponse(w, http.StatusUnauthorized, constants.FAILURE, constants.MissingAPIkey, nil)

			return
		}

		if key != a.apiKey {
			a.logger.Error(constants.InvalidAPIkey)
			a.writeJSONResponse(w, http.StatusBadRequest, constants.FAILURE, constants.InvalidAPIkey, nil)

			return
		}

		h.ServeHTTP(w, r)
	})

	return hf
}

func (a *apiServer) writeJSONResponse(w http.ResponseWriter, status int, respType, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := entities.APIResponse{
		Code:    status,
		Message: message,
		Type:    respType,
	}

	if data != nil {
		res.Data = data
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		a.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getStatusCode(err error) int {
	switch err {
	case constants.ErrInvalidPromoCodeLength:
		return http.StatusUnprocessableEntity
	case constants.ErrInvalidPromoCode:
		return http.StatusBadRequest
	case constants.ErrProductNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
