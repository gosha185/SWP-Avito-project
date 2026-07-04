# Sprint Review Summary – Week 5

**Date:** 4 July 2026

**Attendees:**
- **Customer representative:** Industry Expert
- **Team members:** Georgii, Ekaterina, Stepan, Ivan, Leilia

**Meeting purpose:** Sprint Review – demonstration of Sprint 3 results, User Acceptance Testing (UAT), and feedback on the team's development process.

---

## Sprint Goal Reviewed

**Sprint Goal:** Deliver MVP v2, document architecture (static, dynamic, deployment views), create 3 ADRs, and document development process with git workflow diagram.

**Status:** Partially completed.

---

## Delivered MVP v2 Increment

### Product Changes
- Fixed HTTP 500 → HTTP 400 error handling (#177)
- Added TTL workers for expired point batches (#9)
- Added TTL workers for expired holds (#32)
- Added points expiring in given time (#8)
- Added view all current held points (#23)
- Separated balance endpoint (balance only + expiring points separately)

### Architecture Documentation
- **Static View:** Component diagram showing layered architecture (HTTP → Service → Storage → PostgreSQL)
- **Dynamic View:** Sequence diagram for hold/confirm/cancel flows
- **Deployment View:** Deployment diagram showing Go API service, PostgreSQL, and Docker Compose

### ADRs (Architecture Decision Records)
- ADR-001: Go as backend language
- ADR-002: PostgreSQL as primary database
- ADR-003: Idempotency via external_key

### Development Process Documentation
- Created `docs/development-process.md` with:
  - Git workflow (Mermaid gitGraph diagram)
  - Issue workflow
  - Definition of Done
  - Code Review Process
  - CI and Quality Gates
  - Database and Migrations
  - Testing Strategy
  - CI/CD and Deployment

---

## UAT Results Summary

| UAT Scenario | Status |
|--------------|--------|
| UAT-001: View balance | ✅ Passed |
| UAT-002: View holds | ✅ Passed |
| UAT-003: Award points | ✅ Passed |
| UAT-004: Place hold | ✅ Passed |
| UAT-005: Confirm hold | ✅ Passed |
| UAT-006: Cancel hold | ✅ Passed |
| UAT-007: Cancel already cancelled hold | ✅ Fixed (#177) |

**Note:** New UAT scenarios for MVP v2 were added and executed successfully.

---

## Architecture Evidence Discussed

1. **Component Diagram** – Layered architecture with HTTP, Service, Storage, PostgreSQL, and Worker
2. **Sequence Diagram** – Customer suggested splitting into smaller, scenario-specific diagrams (e.g., balance flow, hold flow with alt blocks for confirm/cancel/timeout)
3. **Deployment Diagram** – Go API service, PostgreSQL, Docker Compose, university network deployment
4. **ADRs** – Discussed with customer; received positive feedback on ADR-002 (background workers as separate server)

**Customer feedback on architecture:**
- Sequence diagrams should be smaller and scenario-specific
- Use alt blocks for hold/confirm/cancel/timeout scenarios
- ADR-002 (separate server for background workers) was well received

---

## Customer Feedback and Approvals

### Positive Feedback
- Good progress on MVP v2
- Architecture documentation is on the right track
- Workers demonstrated working correctly
- Development process documentation is comprehensive

### Suggestions for Improvement
- **Sequence diagrams:** Split into smaller, scenario-specific diagrams (balance flow, hold flow with alt blocks)
- **DoD/DoR templates:** Add checklists to GitHub Issue templates so every new issue automatically includes them
- **Testing:** Use unit tests instead of manual text-based testing for workers
- **Transaction history:** Add time filters (between two dates) as a future improvement
- **Optimisation:** If optimising, measure and describe current vs improved performance

### Approvals
- The customer approved the MVP v2 increment
- The customer approved the architecture direction

---

## Action Points

| Task | Owner |
|------|-------|
| Split sequence diagrams into scenario-specific diagrams | Georgii |
| Add DoD/DoR checklists to GitHub Issue templates | Leilia |
| Fix HTTP 500 issue in confirm/cancel flows | Stepan |
| Use unit tests for worker testing instead of manual | Ivan |
| Add time filters for transaction history (future sprint) | Team |

---

## Risks and Feedback

- No major risks identified.
- The customer appreciated the overall progress and architecture documentation.
- Minor issues with PR review process (reviews approved but still blocked) — needs investigation.
- The customer recommended continuing to use Daily Scrums.

---

## Next Steps

1. Complete remaining tasks (#203 — screenshots, final files)
2. Fix HTTP 500 issue in confirm/cancel flows
3. Split sequence diagrams into smaller scenario-specific diagrams
4. Create final Week 5 reports (reflection, retrospective, llm-report)
5. Create SemVer release for MVP v2
6. Submit Assignment 5 PDF to Moodle
