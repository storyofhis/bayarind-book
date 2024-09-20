# Use the official Golang image as the base image
FROM golang:1.21.2-alpine3.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Important:
#   Because this is a CGO enabled package, you are required to set it as 1.
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

# Copy the source code
COPY . .

# Build the application
RUN go build -o bayarind-book ./cmd/main.go

# Use a smaller base image for the final stage
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bayarind-book /bayarind-book

# Set the entry point for the application
ENTRYPOINT ["./bayarind-book"]