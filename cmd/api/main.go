package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sunimalherath/orderfoodonline/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	api := server.NewAPIServer(server.WithLogger(logger))

	httpHandler := api.RegisterRoutes()

	server := createHTTPServer(httpHandler)

	err := server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("http server error: %s", err.Error()))

		os.Exit(1)
	}
}

func createHTTPServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8090",
		Handler: h,
	}
}
