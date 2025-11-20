# Dockerfile для backend
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./  
RUN go mod download

COPY . .
RUN go build -o /shorty ./cmd/server

# Финальный образ
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /shorty /shorty
ENTRYPOINT ["/shorty"]