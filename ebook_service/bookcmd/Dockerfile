# First stage: Build the Go application
FROM golang:alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go source files
COPY *.go ./

# Install required Go packages (ensure dependencies are installed)
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux

# Build the Go application
RUN go build -o /bookcmd .

# Second stage: Copy the built binary to a clean image
FROM alpine:latest

# Install necessary dependencies to run the binary (e.g., libc)
RUN apk add --no-cache libc6-compat

RUN apk --no-cache add curl

# Copy the binary from the builder stage
COPY --from=builder /bookcmd /bookcmd

# Expose the port the app will run on
EXPOSE 3000

# Command to run the executable
CMD ["/bookcmd"]

