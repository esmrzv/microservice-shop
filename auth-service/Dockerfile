FROM golang:1.26 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth-service ./cmd/api


FROM debian:bookworm-slim
COPY --from=builder /app/auth-service /auth-service
CMD ["/auth-service"]