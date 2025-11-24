// Package adapters: contains all the interfaces
package adapters

import (
	"context"

	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type ProductsRepo interface {
	GetProducts(ctx context.Context) []entities.Product
	GetProductByID(ctx context.Context, productID int64) (*entities.Product, error)
}
