# ============================================================================
# Multi-stage Dockerfile for Bidding System SDK
# ============================================================================

# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/bidding-server ./cmd/server

# ---------- Stage 2: Runtime ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary and migration files
COPY --from=builder /bin/bidding-server .
COPY migrations/ ./migrations/

# Switch to non-root
USER appuser

ENV SERVER_PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./bidding-server"]
