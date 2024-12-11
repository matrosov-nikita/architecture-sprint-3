# Указываем базовый образ с Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main cmd/main.go

# Этап выполнения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .

COPY --from=builder /app/migrations ./migrations

# Указываем порт, который будет использоваться контейнером
EXPOSE 8088

# Команда запуска приложения
CMD ["./main"]
