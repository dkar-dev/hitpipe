FROM golang:1.25-alpine AS builder
WORKDIR /app

# --- Оптимизация кэширования зависимостей ---
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/server/main.go

#FROM gcr.io/distroless/base-debian12
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

COPY ./config/.env ./config/.env
COPY ./config/config.yaml ./config/config.yaml

EXPOSE 8080
CMD ["./main"]