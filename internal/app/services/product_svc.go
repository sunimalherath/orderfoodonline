// Package services: performs the business logic and calls the repository layer to get data.
package services

import (
	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type productSvc struct {
	prodRepo adapters.ProductsRepo
}

func NewProductService(prodRepo adapters.ProductsRepo) adapters.ProductService {
	return &productSvc{
		prodRepo: prodRepo,
	}
}

func (p *productSvc) ListProducts() ([]entities.Product, error) {
	products := p.prodRepo.GetProducts()

	if len(products) == 0 {
		return nil, constants.ErrNoProductsAvailable
	}

	return products, nil
}

func (p *productSvc) FindProductByID(productID int64) (*entities.Product, error) {
	return p.prodRepo.GetProductByID(productID)
}
