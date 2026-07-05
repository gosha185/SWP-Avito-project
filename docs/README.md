# Avito Bonus Service — Documentation

Welcome to the documentation for the Avito Bonus Points Service.

---

## Overview

Bonus Service is a backend service for a bonus points system. It implements earning, holding, and spending bonus points with TTL support, a two-phase spending model (hold / confirm), and full operation auditing via an immutable ledger.

The service is written in Go and uses PostgreSQL as the primary data storage.

---

## Quick Links

- [GitHub Repository](https://github.com/gosha185/SWP-Avito-project)
- [Deployed Service](http://10.93.26.189:8080/)
- [Swagger UI](http://10.93.26.189:8080/)

---

## Documentation

### Process and Project Management

- [Development Process](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/development-process.md) — how the team works
- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md) — completion criteria
- [Roadmap](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md) — sprint-by-sprint plan
- [User Stories](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-stories.md) — feature requirements

### Quality and Testing

- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md) — measurable quality goals
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md) — automated quality checks
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md) — customer scenarios
- [Testing](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/testing.md) — test strategy and coverage

### Architecture

- [Architecture Overview](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/README.md) — system structure
- [Static View](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/static-view/) — component diagram
- [Dynamic View](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/dynamic-view/) — sequence diagrams
- [Deployment View](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/deployment-view/) — deployment structure
- [ADRs](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/adr/) — Architecture Decision Records

### Technical Reference

- [Migrations](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/migrations.md) — database schema changes
- [Storage](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/storage.md) — data access layer

---

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development without Docker)

### Local Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/gosha185/SWP-Avito-project.git
   cd SWP-Avito-project
   ```

2. Copy environment configuration:

   ```bash
   cp .env.example .env
   ```

3. Start the service:

   ```bash
   cd bonus_service
   docker compose up --build -d
   ```

4. Access Swagger UI at `http://localhost:8080/`

---

## Project Status

| Milestone | Status | Date |
|-----------|--------|------|
| Sprint 1 — MVP v1 | ✅ Completed | 15–21 June |
| Sprint 2 — Quality and Testing | ✅ Completed | 22–28 June |
| Sprint 3 — MVP v2 | ✅ Completed | 29 June – 5 July |
| Sprint 4 — Post-MVP Improvements | 📅 Planned | 6–12 July |

---

## License

This project is licensed under the [MIT License](https://github.com/gosha185/SWP-Avito-project/blob/main/LICENSE).
