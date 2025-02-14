# First stage: Build the Go application
FROM golang:1.20-alpine AS builder

# Install dependencies for building (if any additional build tools are needed)
RUN apk update && apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Run go mod tidy to fetch dependencies
RUN go mod tidy

# Copy all Go files into the container
COPY *.go ./

RUN go get github.com/confluentinc/confluent-kafka-go/kafka \
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/gorilla/mux

# Build the Go binary
RUN go build -o /customerservice .

# Second stage: Create the final minimal image
FROM alpine:latest

# Install necessary dependencies to run the app
RUN apk update && apk add --no-cache libc6-compat

# Copy the compiled Go binary from the builder image
COPY --from=builder /customerservice .

# Expose the port the app will run on
EXPOSE 3000

# Command to run the Go application
CMD ["/customerservice"]

