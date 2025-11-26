FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api


FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/bin/api ./api

COPY internal/config/data ./internal/config/data

USER nonroot:nonroot

CMD ["./api"]
