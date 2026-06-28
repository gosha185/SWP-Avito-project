# SWP-Avito-project
## Run with Docker

The MVP v1.1.0 service (Go API + PostgreSQL + self-hosted Swagger UI)  runs with a single command.

    cd bonus_service
    cp .env.example .env          # then set POSTGRES_PASSWORD and other variables if needed
    docker compose up --build -d

- Swagger UI: http://localhost:8080/
- API base path: http://localhost:8080/v1
- Health check: http://localhost:8080/v1/healthcheck

Deployed instance (university network): http://10.93.26.189:8080/

