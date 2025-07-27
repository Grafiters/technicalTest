FROM golang:1.20 AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH="amd64" \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure all dependencies are fetched
RUN go mod tidy

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

RUN go build -o api ./cmd/api/main.go
RUN go build -o task-list ./task/cmd/database.go

FROM alpine:3.13.6

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /build/db ./db
COPY --from=builder /build/api ./
COPY --from=builder /build/task-list ./