# BookShelf REST API

![Go Version](https://img.shields.io/badge/go-1.20%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)
[![Swagger Docs](https://img.shields.io/badge/swagger-docs-brightgreen)](http://localhost:8080/swagger/index.html)

BookShelf - это полнофункциональный RESTful API для управления каталогом книг с авторизацией, избранными подборками и админ-панелью. Приложение реализовано на Go с использованием современных практик разработки и архитектурных паттернов.

## Особенности

- 🛡️ **JWT-аутентификация** с ролевой моделью (user/admin)
- 📚 **CRUD операции** для управления книгами
- ❤️ **Система избранного** для пользователей
- 📊 **Пагинация и фильтрация** книг по жанрам
- 📝 **Swagger документация** для всех эндпоинтов
- 🚀 **Высокая производительность** благодаря кэшированию Redis
- 🐳 **Docker-контейнеризация** для простого развертывания
- ✅ **Полное покрытие тестами** (unit tests с использованием Testify)

## Технологический стек

- **Язык**: Go 1.20+
- **Фреймворк**: Chi Router
- **База данных**: PostgreSQL
- **Кэширование**: Redis
- **Документация**: Swagger (Swaggo)
- **Контейнеризация**: Docker, Docker Compose
- **Аутентификация**: JWT
- **Тестирование**: Testify

## Архитектура

```
/book-catalog
├── .env                  # Environment variables
├── docker-compose.yml    # Docker configuration
├── Dockerfile            # Go app build
├── cmd/
│   └── main.go           # Entry point
├── internal/
│   ├── models/           # DB models (GORM)
│   ├── handlers/         # HTTP handlers
│   ├── repository/       # DB operations
│   ├── service/          # Business logic
│   ├── config/           # App configuration
│   └── middleware/       # Auth and admin middleware
├── pkg/                  # Utility packages
│   ├── utils/            # Response and JWT helpers
│   └── cache/            # Redis cache implementation
└── docs/                 # Generated Swagger docs
```

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.20+ (опционально, если запускать без Docker)

### Запуск без Docker

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL и Redis в Docker:
```bash
docker pull postgres
docker run --name bookshelf-container \
-e POSTGRES_DB=bookshelf_db \
-e POSTGRES_PASSWORD='yourpassword' \
-p 5433:5432 \
-d postgres
  
docker pull redis
docker run --name redis-container \  
-p 6379:6379 \
-d redis
```
3. Создайте .env по образцу:
```bash
PORT="8080"
DSN="host=localhost user=postgres password=yourpassword dbname=bookshelf_db port=5433 sslmode=disable"
JWT_SECRET="JWT_SECRET"
REDIS_URL="redis://localhost:6379"
```

4. Запустите приложение:
```bash
go run cmd/main.go
```

## API Endpoints

### Аутентификация

| Метод | Эндпоинт           | Описание               | Доступ    |
|-------|--------------------|------------------------|-----------|
| POST  | /auth/register     | Регистрация пользователя | Public    |
| POST  | /auth/login        | Вход в систему         | Public    |

### Пользователи

| Метод | Эндпоинт           | Описание                     | Доступ    |
|-------|--------------------|------------------------------|-----------|
| GET   | /users/me          | Получить текущего пользователя | User      |
| GET   | /users/{id}        | Получить пользователя по ID  | User      |
| GET   | /users             | Получить всех пользователей  | Admin     |
| PUT   | /users/{id}/role   | Изменить роль пользователя   | Admin     |
| DELETE| /users/{id}        | Удалить пользователя         | Admin     |

### Книги

| Метод | Эндпоинт       | Описание                     | Доступ    |
|-------|----------------|------------------------------|-----------|
| GET   | /books         | Получить книги с фильтрацией | Public    |
| GET   | /books/{id}    | Получить книгу по ID         | Public    |
| GET   | /books/genres  | Получить все жанры           | Public    |
| POST  | /books         | Создать книгу                | Admin     |
| PUT   | /books/{id}    | Обновить книгу               | Admin     |
| DELETE| /books/{id}    | Удалить книгу                | Admin     |

### Избранное

| Метод | Эндпоинт              | Описание                      | Доступ    |
|-------|-----------------------|-------------------------------|-----------|
| GET   | /favourites           | Получить избранные книги      | User      |
| POST  | /favourites/{bookID}  | Добавить книгу в избранное    | User      |
| DELETE| /favourites/{bookID}  | Удалить книгу из избранного   | User      |

## Примеры запросов

### Регистрация пользователя
```bash
curl -X POST "http://localhost:8080/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"new_user", "password":"strong_password"}'
```

### Аутентификация
```bash
curl -X POST "http://localhost:8080/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"new_user", "password":"strong_password"}'
```

### Добавление книги в избранное
```bash
curl -X POST "http://localhost:8080/favourites/1" \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Получение книг с пагинацией
```bash
curl -X GET "http://localhost:8080/books?page=2&limit=10"
```

## Документация API

Полная документация API доступна через Swagger UI после запуска приложения:

```
http://localhost:8080/swagger/index.html
```

![Swagger UI](https://raw.githubusercontent.com/swagger-api/swagger-ui/master/dist/favicon-32x32.png) [Swagger Documentation](http://localhost:8080/swagger/index.html)

## Тестирование

Для запуска тестов выполните:

```bash
go test -v ./...
```

Тесты покрывают:
- Хендлеры API
- Middleware аутентификации
- Бизнес-логику сервисов
- Работу с репозиториями

## Развертывание в Production

Для production развертывания рекомендуется:

1. Использовать PostgreSQL и Redis в managed-сервисах (AWS RDS, Elasticache и т.д.)
2. Настроить переменные окружения:
   - `DSN` - строка подключения к PostgreSQL
   - `JWT_SECRET` - секрет для генерации JWT
   - `REDIS_URL` - URL для подключения к Redis
3. Использовать reverse proxy (Nginx) для обработки HTTPS

## Вклад в проект

Приветствуются пул-реквесты! Основные шаги:

1. Форкните репозиторий
2. Создайте ветку для вашей фичи (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add some amazing feature'`)
4. Запушьте ветку (`git push origin feature/amazing-feature`)
5. Откройте пул-реквест
