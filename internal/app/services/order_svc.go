package services

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/sunimalherath/orderfoodonline/internal/config"
	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type orderSvc struct {
	productSvc adapters.ProductService
	logger     *slog.Logger
}

type OrderSvcOptions func(*orderSvc)

func WithLogger(logger *slog.Logger) OrderSvcOptions {
	return func(o *orderSvc) {
		o.logger = logger
	}
}

func NewOrderSvc(productSvc adapters.ProductService, opts ...OrderSvcOptions) adapters.OrderService {
	odrSvc := &orderSvc{
		productSvc: productSvc,
	}

	for _, opt := range opts {
		opt(odrSvc)
	}

	return odrSvc
}

func (o orderSvc) PlaceAnOrder(ctx context.Context, orderReq entities.OrderReq) (*entities.Order, error) {
	if err := orderReq.Validate(); err != nil {
		return nil, err
	}

	products := []entities.Product{}

	eGroup, errCtx := errgroup.WithContext(ctx)

	eGroup.Go(func() error {
		err := orderReq.ValidateCouponCode()
		if err != nil {
			if !errors.Is(err, constants.ErrEmptyPromoCode) {
				o.logger.Error(err.Error())

				return err
			}

			return nil
		}

		if validateCouponCode(errCtx, orderReq.CouponCode, config.GetCouponFilePaths()) {
			o.logger.Info("valid coupon code")

			return nil
		}

		return constants.ErrInvalidPromoCode
	})

	eGroup.Go(func() error {
		var err error

		products, err = o.getProductsForOrder(errCtx, orderReq.Items)
		if err != nil {
			return err
		}

		return nil
	})

	if err := eGroup.Wait(); err != nil {
		return nil, err
	}

	return &entities.Order{
		ID:       uuid.New().String(),
		Items:    orderReq.Items,
		Products: products,
	}, nil
}

func (o orderSvc) getProductsForOrder(ctx context.Context, items []entities.OrderItem) ([]entities.Product, error) {
	products := []entities.Product{}

	for _, item := range items {
		prodID, err := strconv.Atoi(item.ProductID)
		if err != nil {
			o.logger.Error(err.Error())

			return products, err
		}

		product, err := o.productSvc.FindProductByID(ctx, int64(prodID))
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func validateCouponCode(ctx context.Context, couponCode string, filePaths []string) bool {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan bool, len(filePaths))

	for _, file := range filePaths {
		go verifyCouponCode(ctx, couponCode, file, resultCh)
	}

	matches := 0

	for range filePaths {
		select {
		case matchFound := <-resultCh:
			if matchFound {
				matches++

				if matches >= 2 {
					cancel()

					return true
				}
			}
		case <-ctx.Done():
			return false
		}
	}

	return false
}

func verifyCouponCode(ctx context.Context, couponCode string, filePath string, resultCh chan<- bool) {
	file, err := os.Open(filePath)
	if err != nil {
		resultCh <- false

		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if scanner.Text() == couponCode {
			resultCh <- true

			return
		}
	}

	resultCh <- false
}
