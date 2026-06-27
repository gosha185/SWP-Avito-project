# Testing Documentation

This document describes the testing strategy, critical modules, coverage status, and CI quality gates for the Avito Bonus Points Service.

---

## Critical Modules and Coverage

Critical modules are source files responsible for core user workflows, persistence, business rules, or other behavior where defects would materially affect the product.

| Critical module | Why critical | Required line coverage | Current line coverage | 
|---|---|---|---:|---:|
| `internal/balance/balance.go` | Core balance operations | 30% | — |
| `internal/holds/holds.go` | Hold creation, confirmation, cancellation | 30% | — | 
| `internal/accrual/accrual.go` | Point accrual with TTL and idempotency | 30% | — |
| `internal/ttl/ttl_worker.go` | TTL worker for expired points and holds | 30% | — |
| `internal/ledger/ledger.go` | Immutable audit ledger operations | 30% | — | 
 Evidence: [Coverage run] (UPDATE)
---

## Automated Test Status

| Test type | Scope | Command or CI check | Latest result | Evidence |
|---|---|---|---|---|
| Unit tests | Critical product logic | `go test ./internal/... -v` | — | [CI run] (UPDATE)  |
| Integration tests | API + database interaction | `go test ./tests/integration/... -v` | — | [CI run] (UPDATE) |
| Automated QRTs | QR-001 to QR-005 | `go test -run=TestPerformance ./...` | — | [QRT report] (UPDATE) |

---

## CI and QA Check Status

| Gate or check | Required for Done? | Latest protected-branch status | Evidence |
|---|---|---|---|
| Linting (`golangci-lint`) | Yes | — | [CI run]  (UPDATE) |
| Build (`go build`) | Yes | — | [CI run] (UPDATE) |
| Unit tests | Yes | — | [CI run]  (UPDATE)|
| Integration tests | Yes | — | [CI run] (UPDATE) |
| Coverage report | Yes | — | [Coverage run] (UPDATE) |
| Automated QRTs | Yes | — | [QRT report] (UPDATE) |
| Additional QA check (dependency vulnerability scan) | Yes | — | [QA check report] (UPDATE) |

---

## Additional QA Check Rationale

| QA objective or risk | Additional QA check | Scope | Latest result | Evidence | Limitations or follow-up |
|---|---|---|---|---|---|
| Dependencies with known vulnerabilities may expose users or deployments to avoidable risk. | Automated dependency vulnerability scan (`govulncheck`) | Product dependency manifests (`go.mod`, `go.sum`) | — | [CI run] (UPDATE) | Some vulnerabilities may require manual triage or delayed upstream fixes. |

---

## Manual Evidence That Does Not Count as QRT

| Evidence | Scope | Result | Follow-up PBI or issue |
|---|---|---|---|
| Customer UAT observation | End-user scenarios (balance, holds, accrual) | — | [UAT results]  (UPDATE) |

---

## Quality Gates That Remain Active

All Assignment 4 quality gates are maintained product assets and will continue to apply during later project work:

- ✅ Linting (`golangci-lint`)
- ✅ Build check (`go build`)
- ✅ Unit tests
- ✅ Integration tests
- ✅ Automated QRTs (QR-001 to QR-005)
- ✅ Coverage reporting
- ✅ Additional QA check (dependency vulnerability scan)

When a product change makes a check obsolete, the team will replace it with an equivalent or stronger check and document the reason.
