# Build stage
FROM golang:alpine AS builder

# Install git, gcc, musl-dev, and sqlite-dev for CGO and SQLite support
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Enable CGO and set the working directory
ENV CGO_ENABLED=1
WORKDIR /go/src/app

# Copy the Go source files
COPY . .

# Download dependencies
RUN go mod download

# Build the application with CGO enabled
RUN go build -o /go/bin/app ./cmd/myapp

# Final stage
FROM alpine:latest

# Install ca-certificates and SQLite runtime library
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the binary from the builder stage
COPY --from=builder /go/bin/app /app

# Set the entry point
ENTRYPOINT ["/app"]

# Expose port 8080
EXPOSE 8080
