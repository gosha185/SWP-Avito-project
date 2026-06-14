# SWP-Avito-project
## Run with Docker (MVP v0)

The MVP v0 service (Go API + PostgreSQL + self-hosted Swagger UI) lives in
`itmo_ledger/` and runs with a single command.

    cd itmo_ledger
    cp .env.example .env          # then set POSTGRES_PASSWORD in .env
    docker compose up --build -d

- Swagger UI: http://localhost:8080/
- API base path: http://localhost:8080/v1
- Health check: http://localhost:8080/v1/healthcheck

Deployed instance (university network): http://10.93.26.189:8080/

See `reports/week2/mvp-v0-report.md` for the MVP v0 report and smoke-check
scenario, and `reports/week2/README.md` for the Week 2 submission index.
