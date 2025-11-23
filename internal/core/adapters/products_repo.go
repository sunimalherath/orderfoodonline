// Package adapters: contains all the interfaces
package adapters

import "github.com/sunimalherath/orderfoodonline/internal/core/entities"

type ProductsRepo interface {
	GetProducts() []entities.Product
	GetProductByID(productID int64) (*entities.Product, error)
}
