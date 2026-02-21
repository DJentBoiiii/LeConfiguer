FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o indexing-service ./cmd/main.go

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /app/indexing-service .
EXPOSE 8082
CMD ["./indexing-service"]