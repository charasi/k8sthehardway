# Stage 1: Build stage
FROM golang:alpine AS build

# Install required dependencies for Go build
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Run go mod tidy to fetch dependencies
RUN go mod tidy

# Copy only the webbff.go file (instead of the entire directory)
COPY webbff.go ./

# Install Go dependencies
RUN go get github.com/go-sql-driver/mysql \
    && go get -u github.com/gorilla/mux \
    && go get -u github.com/shopspring/decimal \
    && go get github.com/golang-jwt/jwt

# Build the Go application
RUN go build -o /webbff

# Stage 2: Final stage (using Alpine)
FROM alpine:latest

# Install necessary runtime dependencies (musl libc, etc.)
RUN apk add --no-cache libgcc

RUN apk --no-cache add curl

# Copy the Go binary from the build stage
COPY --from=build /webbff /webbff

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD ["/webbff"]

