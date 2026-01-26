# ========================
# Build stage
# ========================
FROM golang:1.24-alpine AS builder


WORKDIR /app

# Dependency untuk go mod & build
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev \
    && update-ca-certificates

# Copy go mod files dulu biar layer caching kepake
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary (static, kecil, cocok buat alpine final image)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# ========================
# Final stage
# ========================
FROM alpine:latest

WORKDIR /app

# Runtime cert + timezone
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
