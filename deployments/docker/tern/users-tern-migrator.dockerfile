FROM golang:1.23-alpine AS builder

RUN go install github.com/jackc/tern/v2@latest

COPY ../../../scripts/tern/users_migrations /migrations

ENTRYPOINT ["tern"]