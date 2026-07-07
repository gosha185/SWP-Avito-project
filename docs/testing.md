# Testing Documentation

This document describes the testing strategy, critical modules, coverage status, and CI quality gates for the Avito
Bonus Points Service.

---

## Critical Modules and Coverage

Critical modules are source files responsible for core user workflows, persistence, business rules, or other behavior
where defects would materially affect the product.

| Critical module | Why critical | Required line coverage | Current line coverage | Evidence |
|---------------------------------------------|-------------------------------------------|------------------------|----------------------:|---------------------------------------------------------------------------------------:|
| `bonus-service/internal/storage/balance.go` | Core balance operations | 30% | 57.7% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| `bonus-service/internal/storage/hold.go` | Hold creation, confirmation, cancellation | 30% | 48% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| `bonus-service/internal/storage/batch.go` | Batch management with TTL | 30% | 47.2% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| `bonus-service/internal/storage/ledger.go` | Immutable audit ledger operations | 30% | 55.2% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| `bonus-service/internal/storage/hold_batch.go` | Hold-batch mapping | 30% | 86% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |

**Total Coverage: 39.8%**

---

## Automated Test Status

| Test type | Scope | Command or CI check | Latest result | Evidence |
|-------------------|----------------------------|--------------------------------------|---------------|--------------------------------------------------------------------------------------|
| Unit & Integration tests | Critical product logic | `go test ./internal/storage -v` | ✅ Passing | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |


---

## CI and QA Check Status

| Gate or check | Required for Done? | Latest protected-branch status | Evidence |
|-----------------------------------------------------|--------------------|--------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| Build (`go build`) | Yes | ✅ Passing | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| Unit tests | Yes | ✅ Passing | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| Integration tests | Yes | ✅ Passing | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106041) |
| Coverage report | Yes | 39.8% | [Coverage run](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/coverage.jpg) |
| Automated QRTs | Yes | ✅ Passing | [QRT report](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106045) |
| Additional QA check (dependency vulnerability scan) | Yes | ✅ Passing | [QA check report](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/branch-protection.jpg) |

---

## Additional QA Check Rationale

| QA objective or risk | Additional QA check | Scope | Latest result | Evidence | Limitations or follow-up |
|--------------------------------------------------------------------------------------------|---------------------------------------------------------|---------------------------------------------------|---------------|----------------------------------------------------------------------------------|---------------------------------------------------------------------------|
| Dependencies with known vulnerabilities may expose users or deployments to avoidable risk. | Automated dependency vulnerability scan (`govulncheck`) | Product dependency manifests (`go.mod`, `go.sum`) | ✅ Passing | [CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/28748106057) | Some vulnerabilities may require manual triage or delayed upstream fixes. |

---

## Manual Evidence That Does Not Count as QRT

| Evidence | Scope | Result | Follow-up PBI or issue |
|--------------------------|----------------------------------------------|------------|--------------------------------------------------------------------------------------------------------------------|
| Customer UAT observation | End-user scenarios (balance, holds, accrual) | 6/7 passed | [UAT](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md#user-acceptance-tests) |

---

## Quality Gates That Remain Active

All Assignment 4 quality gates are maintained product assets and will continue to apply during later project work:

- ✅ Build check (`go build`)
- ✅ Unit tests
- ✅ Integration tests
- ✅ Coverage reporting (currently 39.8%)
- ✅ Additional QA check (dependency vulnerability scan)

When a product change makes a check obsolete, the team will replace it with an equivalent or stronger check and document
the reason.
