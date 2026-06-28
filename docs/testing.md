# Testing Documentation

This document describes the testing strategy, critical modules, coverage status, and CI quality gates for the Avito
Bonus Points Service.

---

## Critical Modules and Coverage

Critical modules are source files responsible for core user workflows, persistence, business rules, or other behavior
where defects would materially affect the product.

| Critical module                             | Why critical                              | Required line coverage | Current line coverage |                                                                               Evidence |
|---------------------------------------------|-------------------------------------------|------------------------|----------------------:|---------------------------------------------------------------------------------------:|
| `bonus_service/internal/balance/balance.go` | Core balance operations                   | 30%                    |                 87.5% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |
| `bonus_service/internal/holds/holds.go`     | Hold creation, confirmation, cancellation | 30%                    |                    0% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |
| `bonus_service/internal/accrual/accrual.go` | Point accrual with TTL and idempotency    | 30%                    |                 81.8% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |
| `bonus_service/internal/ttl/ttl_worker.go`  | TTL worker for expired points and holds   | 30%                    |                    0% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |
| `bonus_service/internal/ledger/ledger.go`   | Immutable audit ledger operations         | 30%                    |                    0% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |

**Total Coverage: 53.2%**

---

## Automated Test Status

| Test type         | Scope                      | Command or CI check                  | Latest result | Evidence                                                                             |
|-------------------|----------------------------|--------------------------------------|---------------|--------------------------------------------------------------------------------------|
| Unit tests        | Critical product logic     | `go test ./internal/... -v`          | ✅ Passing     | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)     |
| Integration tests | API + database interaction | `go test ./tests/integration/... -v` | ✅ Passing     | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)     |
| Automated QRTs    | QR-001 to QR-005           | `go test -run=TestPerformance ./...` | ✅ Passing     | [QRT report](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) |

---

## CI and QA Check Status

| Gate or check                                       | Required for Done? | Latest protected-branch status | Evidence                                                                                                                                |
|-----------------------------------------------------|--------------------|--------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| Linting (`golangci-lint`)                           | Yes                | ✅ Passing                      | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)                                                        |
| Build (`go build`)                                  | Yes                | ✅ Passing                      | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)                                                        |
| Unit tests                                          | Yes                | ✅ Passing                      | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)                                                        |
| Integration tests                                   | Yes                | ✅ Passing                      | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)                                                        |
| Coverage report                                     | Yes                | 53.2%                          | [Coverage run](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/coverage.jpg)             |
| Automated QRTs                                      | Yes                | ✅ Passing                      | [QRT report](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352)                                                    |
| Additional QA check (dependency vulnerability scan) | Yes                | ✅ Passing                      | [QA check report](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/branch-protection.jpg) |

---

## Additional QA Check Rationale

| QA objective or risk                                                                       | Additional QA check                                     | Scope                                             | Latest result | Evidence                                                                         | Limitations or follow-up                                                  |
|--------------------------------------------------------------------------------------------|---------------------------------------------------------|---------------------------------------------------|---------------|----------------------------------------------------------------------------------|---------------------------------------------------------------------------|
| Dependencies with known vulnerabilities may expose users or deployments to avoidable risk. | Automated dependency vulnerability scan (`govulncheck`) | Product dependency manifests (`go.mod`, `go.sum`) | ✅ Passing     | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28331026352) | Some vulnerabilities may require manual triage or delayed upstream fixes. |

---

## Manual Evidence That Does Not Count as QRT

| Evidence                 | Scope                                        | Result     | Follow-up PBI or issue                                                                                             |
|--------------------------|----------------------------------------------|------------|--------------------------------------------------------------------------------------------------------------------|
| Customer UAT observation | End-user scenarios (balance, holds, accrual) | 6/7 passed | [UAT](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md#user-acceptance-tests) |

---

## Quality Gates That Remain Active

All Assignment 4 quality gates are maintained product assets and will continue to apply during later project work:

- ✅ Linting (`golangci-lint`)
- ✅ Build check (`go build`)
- ✅ Unit tests
- ✅ Integration tests
- ✅ Automated QRTs (QR-001 to QR-005)
- ✅ Coverage reporting (currently 53.2%)
- ✅ Additional QA check (dependency vulnerability scan)

When a product change makes a check obsolete, the team will replace it with an equivalent or stronger check and document
the reason.
