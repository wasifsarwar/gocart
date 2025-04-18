FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /product-service ./cmd/main.go

# Use a smaller image for the final container
FROM alpine:latest

# Add necessary dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /product-service .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./product-service"] 