# ==========================
# Stage 1: Builder
# ==========================
FROM golang:1.24.5 AS builder

WORKDIR /build

# Copy go mod/sum
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source
COPY . .

# Install golang-migrate CLI dengan tag postgres
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

# Build binaries
RUN go build -o api ./cmd/api/main.go
RUN go build -o worker ./cmd/workers/main.go
RUN go build -o task-list ./task/cmd/database.go

# ==========================
# Stage 2: Final Alpine Image
# ==========================
FROM alpine:3.13.6

# Tambahkan sertifikat CA agar HTTPS bisa digunakan
RUN apk add --no-cache ca-certificates

# Buat direktori kerja
WORKDIR /app

# Salin hasil build dari stage sebelumnya
COPY --from=builder /build/api ./api
COPY --from=builder /build/worker ./worker
COPY --from=builder /build/task-list ./task-list
COPY --from=builder /go/bin/migrate ./migrate

# Tambahkan permission eksekusi
RUN chmod +x /app/api /app/worker /app/task-list /app/migrate

# Default command
CMD ["./api"]
