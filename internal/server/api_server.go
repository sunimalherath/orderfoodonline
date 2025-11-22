package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
)

type apiServer struct {
	logger *slog.Logger
	// add product and order services too
}

type APIServerOptions func(*apiServer)

func WithLogger(logger *slog.Logger) APIServerOptions {
	return func(a *apiServer) {
		a.logger = logger
	}
}

func NewAPIServer(opts ...APIServerOptions) adapters.APIServer {
	apiServer := &apiServer{}

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

	return mux
}

func (a *apiServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("in good health"))
	if err != nil {
		a.logger.Error(err.Error())
	}
}
