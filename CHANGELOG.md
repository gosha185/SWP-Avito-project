# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.0.0] - 2026-06-21

First full release of the bonus service: accrual with expiry, two-phase
redemption with FEFO, balance and history, deployed behind a reverse proxy
with interactive API documentation.

### Added

**Storage and schema**
- Ledger-based schema: `balances`, `batches`, `holds`, `hold_batches`, `ledger`
- Repositories for balances, batches, holds, hold-batches and the ledger
- Pessimistic locking (`SELECT FOR UPDATE`) on balance and batch reads
- Unique `external_key` index on the ledger for idempotency
- FEFO ordering of batches by expiry date
- Background cancellation of expired holds (TTL worker)

**Domain logic**
- `BonusService` that runs each operation inside a single transaction
- Accrual with a configurable TTL, recorded as a new batch
- Two-phase redemption: hold, then confirm or cancel
- FEFO spend: holds draw from the batches that expire soonest
- Idempotency and insufficient-funds checks before any mutation

**HTTP API**
- Endpoints: accrual, hold, confirm, cancel, balance, history
- Request and response DTOs, kept separate from storage models
- Idempotency via the `Idempotency-Key` header
- Input validation with field-level error messages
- Paginated history that does not expose internal keys or metadata
- OpenAPI 3.0 specification and Swagger UI

**Infrastructure**
- Multi-stage Dockerfile and Docker Compose stack (db, api, swagger, proxy)
- Caddy reverse proxy serving the API and Swagger UI on a single port
- Database connection string read from the environment
- Schema applied automatically on first database startup

### Fixed
- Ledger amount was not set on confirm and cancel, which violated the
  `amount > 0` constraint and failed the operation
- Hold accepted a zero expiry window; the duration is now validated
- Server-side errors were not logged, which made failures hard to diagnose

## Contributors

- **Katya Efremova** — data layer and database schema: migrations, repositories,
  indexes, locking, FEFO queries, TTL worker
- **Georgy Sergeev** — service layer: transactions, two-phase redemption, FEFO
  allocation, idempotency and balance checks
- **Stepan Grechinskii** — HTTP layer and delivery: handlers, validation, API contract,
  OpenAPI spec, Docker and reverse-proxy deployment

[1.0.0]: https://github.com/gosha185/SWP-Avito-project/releases/tag/v1.0.0
