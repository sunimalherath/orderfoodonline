package constants

import (
	"fmt"
	"time"
)

// env vairable names
const (
	PORT   = "PORT"
	APIKey = "api_key"
)

// http response types for writing JSON response.
const (
	SUCCESS = "success"
	FAILURE = "faiure"
)

// messages to show in the http response.
const (
	RetrievalFailed  = "could not retrieve products"
	ProductsRcvd     = "products retrieved"
	InvalidProdID    = "invalid product ID"
	ProdNotFound     = "product not found"
	InvalidRequest   = "invalid order detail"
	ValidationFailed = "order validation failed"
	OrderFailed      = "failed place the order"
	OrderPlaced      = "order placed"
	GoodHealth       = "health ok"
)

const CheckHealth = "performing health check"

// auth messages
const (
	MissingAPIkey = "missing API key"
	InvalidAPIkey = "invalid API key"
)

// graceful shtudown messages
const (
	GracefulShutdown = "server shutting down gracefully, press Ctrl+C to force"
	ShutdownComplete = "graceful shutdown complete"
)

// time specific.
const (
	ActiveDuration  time.Duration = 30 * time.Second
	ShutdownTimeout time.Duration = 10 * time.Second
)

// file paths.

const DataDir = "./internal/config/data"

const (
	ProductsFile = "products.json"
	CouponBase1  = "couponbase1"
	CouponBase2  = "couponbase2"
	CouponBase3  = "couponbase3"
)

var (
	ProductsFilePath = fmt.Sprintf("%s/%s", DataDir, ProductsFile)
	CouponFilePath1  = fmt.Sprintf("%s/%s", DataDir, CouponBase1)
	CouponFilePath2  = fmt.Sprintf("%s/%s", DataDir, CouponBase2)
	CouponFilePath3  = fmt.Sprintf("%s/%s", DataDir, CouponBase3)
)
