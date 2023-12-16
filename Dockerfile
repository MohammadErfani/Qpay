# Use a minimal Golang image as the base
FROM golang:1.21-alpine AS builder
# Set the working directory inside the container
WORKDIR /app
# Copy the Go project files into the container
COPY . .
# Build the Go application
RUN go build -o Qpay .
# Use a smaller image for the final build
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/Qpay /app/Qpay
COPY --from=builder /app/database/migration /app/database/migration
