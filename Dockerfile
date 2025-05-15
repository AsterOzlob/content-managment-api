# Используем базовый образ Go
FROM golang:1.24 AS builder

# Установим рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Установка swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерация документации
RUN swag init --dir . --output ./docs

# Собираем бинарник — main.go в корне
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cms-api main.go

# Финальный этап: создание лёгкого образа
FROM gcr.io/distroless/static-debian12

# Копируем собранный бинарник и docs
COPY --from=builder /cms-api /cms-api
COPY --from=builder /app/docs /docs

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/cms-api"]