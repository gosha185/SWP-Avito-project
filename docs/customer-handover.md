# Customer Handover — Bonus Service

**Handover status: Ready for independent use**
This document describes the current handover state of the Bonus Service. It is intended for the customer and any reviewer evaluating the transition.

---

## Product status

The Bonus Service is a fully functional backend service implementing a bonus points system with:

- Bonus accrual with configurable TTL
- Two-phase spending model: hold → confirm / cancel
- FEFO (First Expired, First Out) batch spending strategy
- Idempotency via `external_key`
- Full operation audit via immutable ledger
- Automatic expiration of batches and holds (TTL worker)

The service is deployed and accessible. All core user stories from MVP v1 have been implemented and verified.

---

## How to access the product

### Deployed instance (university network)

The service is currently running and accessible within the university network:

- **Swagger UI / API docs:** <http://10.93.26.189:8080/>
  
This instance runs the latest release from the `main` branch.

### Local deployment

To run the service locally:

```bash
git clone https://github.com/gosha185/SWP-Avito-project.git
cd SWP-Avito-project/bonus_service
cp .env.example .env
docker compose up --build -d
```

The service will be available at <http://localhost:8080/>.

---

## Configuration and environment variables

All configuration is provided via a `.env` file in the `bonus_service/` directory. A template is provided in `.env.example`.

| Variable          | Description                          | Example value                                          |
|-------------------|--------------------------------------|--------------------------------------------------------|
| `POSTGRES_USER`   | PostgreSQL username                  | `admin`                                                |
| `POSTGRES_PASSWORD` | PostgreSQL password                | `admin`                                                |
| `POSTGRES_DB`     | PostgreSQL database name             | `admindb`                                              |
| `DB_DSN`          | Full database connection string      | `postgres://admin:admin@db:5432/admindb?sslmode=disable` |

**Secrets handling:** Never commit `.env` to version control. The file is listed in `.gitignore`. For production deployments, use a secrets manager or CI/CD environment variables rather than a plain `.env` file.

The deployed instance uses its own credentials managed by the team. If the customer takes over the deployment, they must set their own credentials before going live.

---

## Setup and deployment steps

### First-time setup

1. Clone the repository and navigate to `bonus_service/`:

   ```bash
   git clone https://github.com/gosha185/SWP-Avito-project.git
   cd SWP-Avito-project/bonus_service
   ```

2. Copy and configure the environment file:

   ```bash
   cp .env.example .env
   # Edit .env if you need non-default credentials
   ```

3. Start the stack:

   ```bash
   docker compose up --build -d
   ```

   This starts PostgreSQL, applies migrations automatically, and launches the API server.

### Updating to a new release

```bash
git pull origin main
cd bonus_service
docker compose up --build -d
```

### Applying database migrations manually

If running without Docker Compose auto-init:

```bash
migrate -path ./migrations -database "$DB_DSN" up
```

To roll back one migration:

```bash
migrate -path ./migrations -database "$DB_DSN" down 1
```

---

## API usage

The full API is documented via Swagger UI, available at the root URL after startup.

Main endpoints:

| Method | Path                                        | Description                            |
|--------|---------------------------------------------|----------------------------------------|
| POST   | `/accrual`                                  | Add bonus points to a user             |
| POST   | `/hold`                                     | Reserve points for an order            |
| POST   | `/users/{user_id}/holds/{order_id}/confirm` | Finalize a hold (spend points)         |
| POST   | `/users/{user_id}/holds/{order_id}/cancel`  | Cancel a hold (release points)         |
| GET    | `/users/{user_id}/balance`                  | Get available balance                  |
| GET    | `/users/{user_id}/balance/expirations`      | Get expiring points within time window |
| GET    | `/users/{user_id}/holds`                    | Get total held points                  |
| GET    | `/users/{user_id}/history`                  | Get transaction history                |

All mutating requests support idempotency via the `Idempotency-Key` header.

A Postman collection is available in [`api/`](../api/) for manual testing.

---

## Troubleshooting

**Service does not start:**

- Check that Docker and Docker Compose are installed and running.
- Verify `.env` has valid values.
- Check logs: `docker compose logs api`

**Database connection errors:**

- Ensure `DB_DSN` in `.env` matches the PostgreSQL credentials (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`).
- If running without Docker Compose, make sure PostgreSQL is running and accessible.

**Migrations not applied:**

- If starting without Docker Compose auto-init, apply migrations manually (see above).
- Check migration status: `migrate -path ./migrations -database "$DB_DSN" version`

**Health check returns non-200:**

- The service may still be starting. Wait a few seconds and retry.
- Check `docker compose ps` to see container status.

---

## Known limitations and unfinished areas

- **Warning for exceeding limits** — the "warning when user exceeds spending limits" user story is in progress and not yet delivered.
- **No authentication / authorization** — the API is currently unauthenticated. All endpoints are publicly accessible. Adding auth is outside the current MVP scope.
- **Deployed instance is university-network only** — the running instance at `10.93.26.189:8080` is reachable only from within the university network. External access requires a VPN or separate deployment.
- **No automated deployment pipeline** — deployments are performed manually by running `docker compose up --build -d` on the server. A CD pipeline is not yet set up.
- **PostgreSQL persistence** — the Docker Compose stack uses a named volume for PostgreSQL data. Removing the volume (e.g. `docker compose down -v`) will delete all data.

---

## Handover status

**Current level: Ready for independent use**
The customer can deploy and run the service independently using the instructions above. The codebase, documentation, and deployment configuration are all publicly available in the repository.

### What has been transferred / delegated

- Full source code in the public GitHub repository: https://github.com/gosha185/SWP-Avito-project
- All documentation in `docs/` and `reports/`
- Docker Compose deployment configuration
- Database migration files
- OpenAPI specification and Postman collection

### What the team currently retains

- The running deployed instance at `10.93.26.189:8080` (university server managed by the team)
- Server credentials for the deployed instance

### Remaining actions before full transition

- Transfer or decommission the deployed server instance, or provide the customer with access credentials
- Deliver the "warning for exceeding limits" user story (in progress)
- Optional: set up a CD pipeline for automated deployments

The documentation is sufficient for independent deployment and use. Support may still be needed for first-time deployment troubleshooting or if the customer wants to take over the existing server instance.

---

## Related documentation

- [`README.md`](../README.md) — project overview and quick start
- [`docs/development-process.md`](development-process.md) — development workflow and CI/CD
- [`docs/migrations.md`](migrations.md) — database schema details
- [`docs/architecture/README.md`](architecture/README.md) — architecture overview
- [`docs/testing.md`](testing.md) — testing strategy
- [`CHANGELOG.md`](../CHANGELOG.md) — release history
