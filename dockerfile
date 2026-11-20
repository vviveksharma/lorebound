FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o lorebound ./cmd

FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /app/lorebound .

ENV APP_PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
      CMD curl -f "http://127.0.0.1:${APP_PORT}/health" || exit 1

ENTRYPOINT ["/app/lorebound"]