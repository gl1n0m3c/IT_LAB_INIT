# Используем образ Golang как базовый
FROM golang:latest

# Установите рабочую директорию внутри контейнера
WORKDIR /app

# Копируйте go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код проекта в рабочую дирректорию
COPY . .

# Сборка приложения
RUN go build -o ./cmd/main ./cmd/main.go

# Порт, на котором будет работать приложение
EXPOSE 8080

# Меняем рабочу директорию
WORKDIR /app/cmd

# Запускаем приложение приложение
ENTRYPOINT ["./main"]
