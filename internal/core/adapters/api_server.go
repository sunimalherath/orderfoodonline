package adapters

import "net/http"

type APIServer interface {
	RegisterRoutes() http.Handler
	HealthCheck(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
	FindProductByID(w http.ResponseWriter, r *http.Request)
	PlaceAnOrder(w http.ResponseWriter, r *http.Request)
}
