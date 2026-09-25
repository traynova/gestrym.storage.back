# Stage 1: Build binary
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o gestrym-storage main.go

# Stage 2: Minimal runtime
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gestrym-storage .

EXPOSE 8080

CMD ["./gestrym-storage"]
