package entities

import (
	"unicode/utf8"

	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
)

type Order struct {
	ID       string      `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderReq struct {
	Items      []OrderItem `json:"items"`
	CouponCode string      `json:"couponCode"`
}

func (or OrderReq) Validate() error {
	if len(or.Items) == 0 {
		return constants.ErrNoItemsInOrderReqd
	}

	for _, item := range or.Items {
		if item.ProductID == "" {
			return constants.ErrProductItemReqd
		}

		if item.Quantity <= 0 {
			return constants.ErrProductQtyReqd
		}
	}

	return nil
}

func (or OrderReq) ValidateCouponCode() error {
	if len(or.CouponCode) == 0 {
		return constants.ErrEmptyPromoCode
	}

	if codeLength := utf8.RuneCountInString(or.CouponCode); codeLength < 8 || codeLength > 10 {
		return constants.ErrInvalidPromoCodeLength
	}

	return nil
}
