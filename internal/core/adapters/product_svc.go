package adapters

import (
	"context"

	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type ProductService interface {
	ListProducts(ctx context.Context) ([]entities.Product, error)
	FindProductByID(ctx context.Context, productID int64) (*entities.Product, error)
}
