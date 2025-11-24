// Package repositories: contains code that connects to the database and store/provide data to service layer
package repositories

import (
	"context"
	"strconv"
	"sync"

	"github.com/sunimalherath/orderfoodonline/internal/core/adapters"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type productsRepo struct {
	prodCache map[string]entities.Product
	cm        sync.RWMutex
}

func NewProductsRepo(prodCache map[string]entities.Product) adapters.ProductsRepo {
	return &productsRepo{
		prodCache: prodCache,
	}
}

func (p *productsRepo) GetProducts(ctx context.Context) ([]entities.Product, error) {
	products := []entities.Product{}

	if err := ctx.Err(); err != nil {
		return products, err
	}

	p.cm.RLock()

	for _, prod := range p.prodCache {
		products = append(products, prod)
	}

	p.cm.RUnlock()

	return products, nil
}

func (p *productsRepo) GetProductByID(ctx context.Context, productID int64) (*entities.Product, error) {
	strProdID := strconv.FormatInt(productID, 10)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p.cm.RLock()
		prod, found := p.prodCache[strProdID]
		p.cm.RUnlock()

		if !found {
			return nil, constants.ErrProductNotFound
		}

		return &prod, nil
	}
}
