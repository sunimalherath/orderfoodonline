package adapters

import "github.com/sunimalherath/orderfoodonline/internal/core/entities"

type ProductService interface {
	ListProducts() ([]entities.Product, error)
	FindProductByID(productID int64) (*entities.Product, error)
}
