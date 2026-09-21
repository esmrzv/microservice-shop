FROM golang:1.26

WORKDIR /app

COPY product-service/go.mod product-service/go.sum ./product-service/
COPY proto/go.mod ./proto/

WORKDIR /app/product-service

RUN go mod download

WORKDIR /app

COPY product-service/ ./product-service/
COPY proto/ ./proto/

WORKDIR /app/product-service

RUN go build -o product-service ./cmd/api

CMD ["./product-service"]