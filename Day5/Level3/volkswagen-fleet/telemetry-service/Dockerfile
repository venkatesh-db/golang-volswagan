
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o telemetry-service

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/telemetry-service .
EXPOSE 8080
CMD ["./telemetry-service"]
