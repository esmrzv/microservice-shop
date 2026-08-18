FROM golang:1.26-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o product-service ./cmd/api
EXPOSE "8081"
CMD ["./product-service"]
