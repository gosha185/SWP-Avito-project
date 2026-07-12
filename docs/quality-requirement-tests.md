# Quality Requirement Tests

This document defines automated quality requirement tests (QRTs) that verify the measurable quality requirements defined in [quality-requirements.md](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md).

---

## QRT-001: API response time

**Linked quality requirement:** [QR-001](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md#qr-001-api-response-time)

**Verification method:** Automated CI check using performance tests.

**Test data, setup, or environment:** Standard CI environment with test database and mock data.

**Automated command or CI check:** 
```bash
cd bonus_service
go test  ./internal/service
```

**Expected measurable result:** 95% of requests for each endpoint complete within 500ms under 100 concurrent users.

**Evidence link:** [Latest CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/29169873511)

---

## QRT-002: Health check availability

**Linked quality requirement:** [QR-002](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md#qr-002-service-availability)

**Verification method:** Automated CI check that verifies health endpoint.

**Test data, setup, or environment:** Deployed service instance.

**Automated command or CI check:** `curl -f http://localhost:8080/healthcheck`

**Expected measurable result:** HTTP 200 OK response.

**Evidence link:** [Latest CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/29169873511)

---

## QRT-003: Critical module unit coverage

**Linked quality requirement:** [QR-003](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md#qr-003-service-availability)

**Verification method:** Automated CI coverage check.

**Test data, setup, or environment:** Standard CI environment.

**Automated command or CI check:** 
```bash
cd bonus_service
go test ./internal/storage
```

**Expected measurable result:** Each critical module has ≥ 30% line coverage.

**Evidence link:** [Latest coverage report](https://github.com/gosha185/SWP-Avito-project/actions/runs/29169873511)

---

## QRT-004: Data retention and integrity

**Linked quality requirement:** [QR-004](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md#qr-004-data-retention-and-availability)

**Verification method:** Automated CI check that verifies data integrity.

**Test data, setup, or environment:** Standard CI environment.

**Automated command or CI check:** 
```bash
cd bonus_service
go test ./internal/service
```

**Expected measurable result:** All transactions from the last 30 days are returned without data loss or corruption.

**Evidence link:** [Latest CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/29169873511)

---

## QRT-005: Fault tolerance and recovery

**Linked quality requirement:** [QR-005](https://github.com/gosha185/SWP-Avito-project/blob/132-quality-requirements/docs/quality-requirements.md#qr-005-data-retention-and-availability)

**Verification method:** Automated CI check that simulates database failure.

**Test data, setup, or environment:** Standard CI environment with simulated database failure.

**Automated command or CI check:**
```bash
cd bonus_service
go test ./internal/service
```

**Expected measurable result:** The system recovers and completes the operation within 5 seconds without data loss.

**Evidence link:** [Latest CI run](https://github.com/gosha185/SWP-Avito-project/actions/runs/29169873511)
