
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init telecom && go mod tidy
RUN go build -o customer-service

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/customer-service .
EXPOSE 8090
CMD ["./customer-service"]
