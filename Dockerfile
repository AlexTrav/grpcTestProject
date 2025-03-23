# Используем официальный образ Golang
FROM golang:1.23

# Рабочая директория внутри контейнера
WORKDIR /app

# Кэш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё
COPY . .

# Сборка gRPC сервера
RUN go build -o server ./grpc_server

# Открываем порт gRPC
EXPOSE 50051

# Запускаем gRPC сервер
CMD ["/app/server"]
