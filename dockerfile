# stage de build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server ./cmd/api

# stage final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server ./server
COPY .env ./
EXPOSE 3000
ENTRYPOINT ["./server"]