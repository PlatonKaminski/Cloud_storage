FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# Финальный образ — минимальный
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /build/app .

COPY --from=builder /build/web ./web
# Копируем миграции
COPY --from=builder /build/migration ./migration

# Создаём папку для файлов
RUN mkdir -p /app/data/files

EXPOSE 8080

CMD ["./app"]