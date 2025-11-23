// Package constants: stores constants and errors to use internally.
package constants

import "errors"

var (
	ErrProductNotFount     = errors.New("product not fount")
	ErrNoProductsAvailable = errors.New("no products available")
	ErrReadingJSONfile     = errors.New("error reading json file")
	ErrUnmarshallingData   = errors.New("error occurred when unmarshalling data")
)
