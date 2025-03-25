# === Dockerfile for kai_security ===

# 1. Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kai_security cmd/main.go

# 2. Final runtime stage
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/kai_security .

EXPOSE 8080

ENTRYPOINT ["/app/kai_security"]
CMD ["-db=/data/kai_security.db", "-port=8080"]