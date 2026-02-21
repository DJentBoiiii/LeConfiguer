FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o config-storage-service ./cmd/main.go

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /app/config-storage-service .
EXPOSE 8081
CMD ["./config-storage-service"]