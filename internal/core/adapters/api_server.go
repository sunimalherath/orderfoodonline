package adapters

import "net/http"

type APIServer interface {
	RegisterRoutes() http.Handler
	HealthCheck(w http.ResponseWriter, r *http.Request)
}
