package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sunimalherath/orderfoodonline/internal/config"
	"github.com/sunimalherath/orderfoodonline/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.Load()

	api := server.NewAPIServer(server.WithLogger(logger))

	httpHandler := api.RegisterRoutes()

	server := createHTTPServer(httpHandler, cfg.Server.Port)

	err := server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("http server error: %s", err.Error()))

		os.Exit(1)
	}
}

func createHTTPServer(h http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: h,
	}
}
