FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o api-gateway-service ./cmd/main.go

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /app/api-gateway-service .

EXPOSE 8080
CMD ["./api-gateway-service"]