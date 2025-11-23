// Package constants: stores constants and errors to use internally.
package constants

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrNoProductsAvailable = errors.New("no products available")
	ErrReadingJSONfile     = errors.New("error reading json file")
	ErrUnmarshallingData   = errors.New("error occurred when unmarshalling data")
)
