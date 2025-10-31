# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o notification-service

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification-service .
EXPOSE 8082
CMD ["./notification-service"]

