# grpcTestProject - gRPC-сервис на Go с PostgreSQL и Docker

Полноценный пример gRPC-приложения на Go с CRUD-операциями пользователей, логированием, хранением данных в PostgreSQL и развёртыванием через Docker.

---

## 🔧 Стек технологий

- **Go (Golang)** — язык программирования
- **gRPC** — взаимодействие между клиентом и сервером
- **Protocol Buffers (proto3)** — описание API
- **PostgreSQL** — база данных
- **GORM** — ORM-библиотека для Go
- **Docker & Docker Compose** — контейнеризация

---

## 📂 Структура проекта

```
grpcTestProject/
├── grpc_client/             # gRPC клиент (консольный)
│   └── main.go              # CLI для взаимодействия с сервером
├── grpc_server/             # gRPC сервер
│   ├── main.go              # Основная логика сервера
│   └── database.go          # Подключение и миграции БД
├── userpb/                  # Сгенерированные protobuf-файлы
├── user.proto               # Описание gRPC API
├── Dockerfile               # Docker-файл для сборки сервера
├── docker-compose.yml       # Поднимает сервер и PostgreSQL
├── go.mod, go.sum           # Модули и зависимости Go
```

---

## 🚀 Быстрый старт

### 1. Клонирование репозитория
```bash
git clone https://github.com/AlexTrav/grpcTestProject.git
cd grpcTestProject
```

### 2. Генерация кода из .proto

Убедитесь, что установлен `protoc` и плагин для Go:

```bash
# Установить protoc-gen-go и protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавить их в PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Сгенерировать код из user.proto
protoc --go_out=. --go-grpc_out=. user.proto
```

### 3. Сборка и запуск в Docker
```bash
docker-compose up --build
```

> Это поднимет:
> - PostgreSQL (порт 5432)
> - gRPC сервер (порт 50051)

---

## 🖥️ Использование клиента

### Запуск клиента:
```bash
go run ./grpc_client/main.go
```

В консоли появится меню:
```
1. Создать пользователя
2. Получить пользователя по ID
3. Обновить пользователя
4. Удалить пользователя
5. Показать всех пользователей
0. Выход
```

---

## 📦 Пример gRPC-запросов

### Протобуфер (`user.proto`):
```proto
service UserService {
  rpc CreateUser (UserRequest) returns (UserResponse);
  rpc GetUser (UserID) returns (UserResponse);
  rpc UpdateUser (UpdateRequest) returns (UserResponse);
  rpc DeleteUser (UserID) returns (DeleteResponse);
  rpc ListUsers (Empty) returns (UserList);
}
```

---

## ⚙️ Переменные окружения (docker-compose)

```yaml
environment:
  DB_HOST: postgres
  DB_PORT: 5432
  DB_USER: postgres
  DB_PASSWORD: root
  DB_NAME: go_test_db
```

---

## 🧪 Пример подключения к БД (GORM)
```go
dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
  getEnv("DB_HOST", "localhost"),
  getEnv("DB_USER", "postgres"),
  getEnv("DB_PASSWORD", "root"),
  getEnv("DB_NAME", "go_test_db"),
  getEnv("DB_PORT", "5432"),
)
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

---

## 🧹 Остановка и очистка контейнеров
```bash
docker-compose down --volumes --remove-orphans
```

---

## 📌 Примечания

- Временное ограничение (`context.WithTimeout`) в клиенте установлено для предотвращения зависания
- Все логи сервера выводятся в консоль Docker или локального запуска
- Полезно для изучения gRPC + Go + PostgreSQL

---

## ✅ Статус

✔️ Проект полностью работает.  
✔️ Подходит для демонстрации навыков работы с gRPC + Go + PostgreSQL + Docker.  
✔️ Можно использовать как шаблон для других проектов.

---

## 👨‍💻 Автор

> Разработано в рамках тестового задания.  
> GitHub: [AlexTrav](https://github.com/AlexTrav)  
> Telegram: [@alex_codov](https://t.me/alex_codov)

---

💬 Вопросы? Предложения? Добро пожаловать в Issues или Telegram!
