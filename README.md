# Order Food Online API

This is a solution for the "Advance Challenge" of the Oolio group's back-end challenge. 

The description of the challenge can be found in [here.](https://github.com/oolio-group/kart-challenge/tree/advanced-challenge/backend-challenge)
API Schema for the solution can be found in [here.](https://github.com/sunimalherath/orderfoodonline/blob/main/docs/openapi.yaml) 
## API Endpoints

`GET /health`                 # Perfom health check.
`GET /product`                # List all the products.
`GET /product/{productId}`    # List product details for the provided `productId`.
`POST /order`                 # Place an order.


## Prerequisites

- Golang v1.24.3 or higher.

## Quick Start

> [!NOTE]
> Replace the couponbase1, couponbase2, and couponbase3 files with their larger couterparts.  

Start the API server:

`make run`
or
`go run cmd/api/main.go`

Start using Docker Compose:

`make docker-up`
or
`docker-compose up --build -d`

`make test`
or
`go test -v ./internal/server`
