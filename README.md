# Bonus Service (SWP Avito Project)

Bonus Service is a backend service for a bonus points system. It implements earning, holding, and spending bonus points with TTL support, a two-phase spending model (hold / confirm), and full operation auditing via an immutable ledger.

The service is written in Go and uses PostgreSQL as the primary data storage.

---

## Features

### Business capabilities

- Accrual of bonus points
- Retrieval of user balance
- Two-phase spending model:
  - hold: reserve points for an order
  - confirm: finalize spending
  - cancel: release reserved points
- FEFO (First Expired, First Out) batch spending strategy
- TTL support for:
  - point batches
  - active holds
- Transaction history via immutable ledger

### System properties

- Idempotency via external_key
- ACID guarantees via PostgreSQL transactions
- Immutable audit log (ledger)
- Indexed data model for performance optimization
- Schema designed for future sharding (composite primary keys)

### Current feature status

**Ready for independent use**
The customer can deploy and operate the service independently. See [`docs/customer-handover.md`](../../docs/customer-handover.md) for full handover details including remaining actions
---

## Quick start (Docker)

### 1. Prepare environment

```bash
cd bonus_service
cp .env.example .env
```

PostgreSQL credentials should be like this:

```POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=admindb
```

Database connection string (used by application) should be like this:

```bash
DB_DSN=postgres://admin:admin@db:5432/admindb?sslmode=disable
```

---

### 2. Run the service

```bash
docker compose up --build -d
```

---

### 3. Access points

- Swagger UI: <http://localhost:8080/>
- Deployed instance (university network): <http://10.93.26.189:8080/>

---

## Architecture

Main modules:

- balance: user balance management
- accrual: bonus accrual logic
- holds: two-phase spending workflow
- ledger: immutable audit log
- ttl worker: expiration processing for batches and holds

---

## Database schema

PostgreSQL is used as the primary database.

Main tables:

- balances: user balance state
- batches: bonus point batches with expiration
- holds: reserved funds for operations
- hold_batches: mapping between holds and batches
- ledger: immutable operation log

See docs/migrations.md for full schema details.

---

## API

The API is defined using OpenAPI 3.0.

Main endpoints:

- POST /accrual
- POST /hold
- POST /users/{user_id}/holds/{order_id}/confirm
- POST /users/{user_id}/holds/{order_id}/cancel
- GET /users/{user_id}/balance
- GET /users/{user_id}/balance/expirations
- GET /users/{user_id}/holds
- GET /balance/{user_id}/history

Swagger UI is available after service startup.

---

## Testing

Run unit and integration tests:

```bash
go test ./internal/... -v
go test ./tests/integration/... -v
```

Performance and QRT tests:

```bash
go test -run=TestPerformance ./... -bench=. -benchtime=10s
```

---

## Quality gates

The project includes CI quality checks:

- linting (golangci-lint)
- build verification
- unit tests
- integration tests
- performance QRT tests
- coverage reporting
- vulnerability scanning (govulncheck)

---

## User stories (MVP)

- view available balance
- view points per order
- admin transaction history
- warning for exceeding limits (in progress)

---

## Configuration

Environment variables:

- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_DB
- DB_DSN

See .env.example for full configuration list.

---

## Docker

The system includes:

- PostgreSQL database
- Go API service
- Automatic database migrations

Persistent storage is configured for PostgreSQL data.

---

## Key implementation details

- FEFO-based batch spending
- TTL worker for automatic cleanup and expiration
- Immutable ledger for auditability
- Idempotency protection at database level

---

## Documentation

- [Development Process](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/development-process.md)
- [Architecture Documentation](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/README.md)
- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md)
- [Roadmap](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md)
- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md)
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md)
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md)
- [Testing Documentation](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/testing.md)
- [Hosted Documentation](https://gosha185.github.io/SWP-Avito-project/)
- [docs/customer-handover.md](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/customer-handover.md)
- [CONTRIBUTING.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CONTRIBUTING.md)
- [AGENTS.md](https://github.com/gosha185/SWP-Avito-project/blob/main/AGENTS.md)
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)
- [docs/development-process.md](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/development-process.md)
