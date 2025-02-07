# Use latest golang as a base image
FROM golang:latest AS builder

# Set directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port being used by application
EXPOSE 8080

# Compile the app
RUN go build -o /receipt-api

# Command to run the executable
CMD ["/receipt-api"]