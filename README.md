# Simple ToDo list in go with htmx

## Technology Stack

- Go
- Docker
- PostgreSQL
- HTMX
- Gin Gonic
- Zap Logger
- pgxPool
- squirrel

## Running the Project

### Using Docker

1. Ensure you have Docker and Docker Compose installed on your machine.
2. Clone the repository:
   ```sh
   git clone https://github.com/netscrawler/simpletodo.git
   cd simpletodo
   ```
3. Build and run the Docker containers:
   ```sh
   docker-compose up --build
   ```
4. The application will be available at `http://localhost:8080`.

### Running Locally without Docker

1. Ensure you have Go and PostgreSQL installed on your machine.
2. Clone the repository:
   ```sh
   git clone https://github.com/netscrawler/simpletodo.git
   cd simpletodo
   ```
3. Set up the PostgreSQL database:
   ```sh
   psql -U postgres -c "CREATE DATABASE htmxtst;"
   psql -U postgres -d htmxtst -f migrations/init.sql/001_createTables.sql
   ```
4. Set the environment variables:
   ```sh
   export DATABASE_HOST=localhost
   export DATABASE_PORT=5432
   export DATABASE_USER=postgres
   export DATABASE_PASSWORD=postgres
   export DATABASE_NAME=htmxtst
   export SERVER_PORT=8080
   export SERVER_HOST=0.0.0.0
   ```
5. Build and run the application:
   ```sh
   go build -o main ./cmd/main.go
   ./main
   ```
6. The application will be available at `http://localhost:8080`.
