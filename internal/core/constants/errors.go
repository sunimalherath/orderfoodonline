// Package constants: stores constants and errors to use internally.
package constants

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrNoProductsAvailable = errors.New("no products available")
	ErrReadingJSONfile     = errors.New("error reading json file")
	ErrUnmarshallingData   = errors.New("error occurred when unmarshalling data")
	ErrWritingResponse     = errors.New("error occurred when writing response")
)

// validation errors
var (
	ErrNoItemsInOrderReqd = errors.New("items required for the order request")
	ErrProductItemReqd    = errors.New("productId cannot be empty")
	ErrProductQtyReqd     = errors.New("product quantity required")

	ErrEmptyPromoCode         = errors.New("promo code cannot be empty")
	ErrInvalidPromoCodeLength = errors.New("invalid promo code length")
	ErrInvalidPromoCode       = errors.New("invalid promo code")
)
