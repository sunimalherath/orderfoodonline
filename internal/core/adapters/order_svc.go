package adapters

import (
	"context"

	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type OrderService interface {
	PlaceAnOrder(ctx context.Context, orderReq entities.OrderReq) (*entities.Order, error)
}
