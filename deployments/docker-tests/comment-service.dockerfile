FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/comment-service ./
COPY internal/comment-service ./internal/comment-service
COPY internal/shared ./internal/shared

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest AS final

COPY --from=builder /app/server /app/server

ENTRYPOINT ["/app/server"]