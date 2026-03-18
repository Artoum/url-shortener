# URL Shortener 🔗

Сервис для сокращения ссылок, написанный на Go. Позволяет создавать короткие URL, устанавливать свои алиасы и выполнять редирект.

## 🚀 Возможности

- **Сокращение ссылок** — превращай длинные URL в короткие
- **Свои алиасы** — задавай запоминающиеся имена (например, `/google`)
- **Автогенерация** — если алиас не указан, создаётся случайный из 6 символов
- **PostgreSQL** — надёжное хранение данных
- **REST API** — простые JSON-эндпоинты
- **Basic Auth** — защита операций записи (POST, DELETE)
- **Структурированные логи** — цветные для разработки, JSON для продакшена
- **Тесты** — юнит-тесты с моками
- **Чистая архитектура** — всё разложено по полочкам

## 🛠️ Технологии

- **Язык:** Go 1.25+
- **Роутер:** [chi](https://github.com/go-chi/chi) — лёгкий и быстрый
- **База данных:** PostgreSQL (драйвер [lib/pq](https://github.com/lib/pq))
- **Конфиги:** [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- **Валидация:** [go-playground/validator](https://github.com/go-playground/validator)
- **Логирование:** slog + кастомный pretty-handler
- **Тестирование:** testify + моки

## 📋 Требования

- Go 1.25 или выше
- PostgreSQL (локально или удалённо)
- Git

## 🚦 Быстрый старт

### 1. Клонируй репозиторий

```bash
git clone https://github.com/artoum/url-shortener.git
cd url-shortener
```

### 2. Настрой PostgreSQL

```bash
# Создай базу данных
createdb -U postgres urlshortener

# или через psql:
psql -U postgres -c "CREATE DATABASE urlshortener;"
```

### 3. Настрой конфиг

Создай файл `config/local.yaml`:

```yaml
env: "local"
database_url: "postgres://postgres:your_password@localhost:5432/urlshortener?sslmode=disable"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 60s
  user: "admin"      # для Basic Auth
  password: "admin"  # для Basic Auth
```

### 4. Установи зависимости и запусти

```bash
go mod download
go run cmd/url-shortener/main.go
```

## 📬 API Эндпоинты

### 🔓 Открытые (без аутентификации)

| Метод | Путь        | Описание                     |
|-------|-------------|------------------------------|
| `GET` | `/{alias}`  | Редирект на оригинальный URL |

### 🔐 Защищенные (Basic Auth)

| Метод   | Путь             | Описание                |
|---------|------------------|-------------------------|
| `POST`  | `/url`           | Создать короткую ссылку |
| `DELETE`| `/url/{alias}`   | Удалить ссылку          |

## 📝 Примеры запросов

### Создать ссылку

```bash
curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d '{"url": "https://google.com"}'
```

**Ответ:**

```json
{
    "status": "OK",
    "alias": "aB3xK9"
}
```

### Создать ссылку со своим алиасом

```bash
curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d '{"url": "https://google.com", "alias": "google"}'
```

### Перейти по короткой ссылке

```
http://localhost:8082/google
```

### Удалить ссылку

```bash
curl -X DELETE http://localhost:8082/url/google -u admin:admin
```

## 🗂️ Структура проекта

```
url-shortener/
├── cmd/
│   └── url-shortener/          # точка входа
├── config/                      # конфиг-файлы
├── internal/
│   ├── config/                  # загрузка конфига
│   ├── http-server/
│   │   ├── handlers/            # хендлеры (save, redirect, delete)
│   │   └── middleware/          # middleware (логгер)
│   ├── lib/                      # утилиты (логи, рандом, ответы)
│   └── storage/                  # работа с БД (postgres)
└── storage/                      # файлы БД (для SQLite, если нужно)
```

## 🧪 Тесты

```bash
# Запустить все тесты
go test ./...

# С покрытием
go test -cover ./...

# Запустить конкретный тест
go test -v ./internal/http-server/handlers/url/save
```

## 📦 Сборка

```bash
# Собрать бинарник
go build -o bin/url-shortener cmd/url-shortener/main.go

# Запустить
./bin/url-shortener
```
