FROM golang:1.26

WORKDIR /app

COPY cart-service/go.mod cart-service/go.sum ./cart-service/
COPY proto/go.mod ./proto/

WORKDIR /app/cart-service

RUN go mod download

WORKDIR /app

COPY cart-service/ ./cart-service/
COPY proto/ ./proto/

WORKDIR /app/cart-service

RUN go build -o cart-service ./cmd/api

CMD ["./cart-service"]