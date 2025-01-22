FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY migrations ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go
RUN go build -o migrateUP ./cmd/migrator/migrateUP.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrateUP .
COPY --from=builder /app/config.env .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000
CMD ["sh", "-c", "sleep 3 && ./migrateUP --storage-path postgres://evans:evans@db_wallet:8001/postgres?sslmode=disable --migrations-path ./migrations && ./main -config ./config.env"]
# ENTRYPOINT ["./migrateUP", "--storage-path", "postgres://evans:evans@db_wallet:8001/postgres?sslmode=disable", "--migrations-path", "./migrations"]
# CMD ["./main", "-config", "./config.env"]