# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go mod init stock_exchange_app && \
    go get github.com/gorilla/mux@v1.8.1 && \
    go get github.com/joho/godotenv@v1.5.1 && \
    go get github.com/go-sql-driver/mysql@v1.8.1

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 2137
CMD ["./server"]
