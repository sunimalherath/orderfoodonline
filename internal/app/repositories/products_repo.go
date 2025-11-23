// Package repositories: contains code that connects to the database and store/provide data to service layer
package repositories

import (
	"strconv"
	"sync"

	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

var cacheMutex sync.RWMutex

type productsRepo struct {
	prodCache map[string]entities.Product
	cm        sync.RWMutex
}

func NewProductsRepo(prodCache map[string]entities.Product) adapters.ProductsRepo {
	return &productsRepo{
		prodCache: prodCache,
	}
}

func (p *productsRepo) GetProducts() []entities.Product {
	products := []entities.Product{}

	p.cm.RLock()

	for _, prod := range p.prodCache {
		products = append(products, prod)
	}

	p.cm.RUnlock()

	return products
}

func (p *productsRepo) GetProductByID(productID int64) (*entities.Product, error) {
	strProdID := strconv.FormatInt(productID, 10)

	p.cm.RLock()
	prod, found := p.prodCache[strProdID]
	p.cm.RUnlock()

	if !found {
		return nil, constants.ErrProductNotFount
	}

	return &prod, nil
}
