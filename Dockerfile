FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk update && apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/deployment ./deployment

# Expose port (adjust if necessary)
EXPOSE 8080

CMD ["./main"]
