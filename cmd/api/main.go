package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunimalherath/orderfoodonline/internal/app/repositories"
	"github.com/sunimalherath/orderfoodonline/internal/app/services"
	"github.com/sunimalherath/orderfoodonline/internal/config"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.Load()

	productCache, err := config.LoadProducts()
	if err != nil {
		logger.Error(err.Error())
	}

	productsRepo := repositories.NewProductsRepo(productCache)

	productSvc := services.NewProductService(productsRepo)
	orderSvc := services.NewOrderSvc(productSvc, services.WithLogger(logger))

	api := server.NewAPIServer(productSvc, orderSvc, server.WithLogger(logger))

	httpHandler := api.RegisterRoutes()

	server := createHTTPServer(httpHandler, cfg.Server.Port)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			logger.Error(fmt.Sprintf("http server error: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Info(constants.GracefulShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(constants.ShutdownComplete)
}

func createHTTPServer(h http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: h,
	}
}
