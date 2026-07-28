FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/api \
    ./cmd/api

FROM alpine:3.24

WORKDIR /app

RUN adduser -D -H appuser

COPY --from=builder --chown=appuser:appuser /app/api /app/api

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/api"]
