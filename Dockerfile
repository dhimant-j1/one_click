# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy database schema and any other necessary files
COPY --from=builder /app/database ./database

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"]
