# Простой список задач на Go с использованием htmx

## Технологический стек

- Go
- Docker
- PostgreSQL
- HTMX
- Gin Gonic
- Zap Logger
- pgxPool
- squirrel

## Запуск проекта

### Использование Docker

1. Убедитесь, что у вас установлены Docker и Docker Compose.
2. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/netscrawler/simpletodo.git
   cd simpletodo
   ```
3. Соберите и запустите контейнеры Docker:
   ```sh
   docker-compose up --build
   ```
4. Приложение будет доступно по адресу `http://localhost:8080`.

### Запуск локально без Docker

1. Убедитесь, что у вас установлены Go и PostgreSQL.
2. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/netscrawler/simpletodo.git
   cd simpletodo
   ```
3. Настройте базу данных PostgreSQL:
   ```sh
   psql -U postgres -c "CREATE DATABASE htmxtst;"
   psql -U postgres -d htmxtst -f migrations/init.sql/001_createTables.sql
   ```
4. Установите переменные окружения:
   ```sh
   export DATABASE_HOST=localhost
   export DATABASE_PORT=5432
   export DATABASE_USER=postgres
   export DATABASE_PASSWORD=postgres
   export DATABASE_NAME=htmxtst
   export SERVER_PORT=8080
   export SERVER_HOST=0.0.0.0
   ```
5. Соберите и запустите приложение:
   ```sh
   go build -o main ./cmd/main.go
   ./main
   ```
6. Приложение будет доступно по адресу `http://localhost:8080`.
