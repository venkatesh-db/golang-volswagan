# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o analytics-service

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/analytics-service .
EXPOSE 8081
CMD ["./analytics-service"]
