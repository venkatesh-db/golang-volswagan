# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o dashboard-service

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dashboard-service .
EXPOSE 8083
CMD ["./dashboard-service"]
