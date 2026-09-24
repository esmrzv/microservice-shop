FROM golang:1.26

WORKDIR /app

COPY product-service/go.mod product-service/go.sum ./product-service/
COPY proto/go.mod ./proto/

WORKDIR /app/product-service

RUN go mod download
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

WORKDIR /app

COPY product-service/ ./product-service/
COPY proto/ ./proto/

WORKDIR /app/product-service

RUN go build -o product-service ./cmd/api

CMD ["./product-service"]