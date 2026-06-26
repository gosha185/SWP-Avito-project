# Definition of Done

A PBI (Product Backlog Item) is considered **Done** when all of the following conditions are met:

---

## 1. Issue-specific Acceptance Criteria

- All acceptance criteria defined in the issue are satisfied and verified.

---

## 2. Code Review

- The work is reviewed and approved by at least **one other team member**.
- The reviewer has confirmed that the implementation meets the acceptance criteria.

---

## 3. CI Quality Gates

- All required CI checks pass for the pull request and the protected default branch:

| Gate | What it checks |
|------|----------------|
| Linting (`golangci-lint`) | Code style and syntax errors |
| Build (`go build`) | Code compiles successfully |
| Unit tests | Critical product logic works |
| Integration tests | Component interactions work |
| Coverage report | Critical modules have ≥ 30% line coverage |
| Automated QRTs | Quality requirements are verified |
| Additional QA check | Dependency vulnerability scan |

---

## 4. Testing Evidence

- Testing evidence is preserved in the PR/MR description, CI logs, or linked documentation.
- The `docs/testing.md` file reflects the current test status for critical modules.

---

## 5. PR/MR Merged

- The issue-linked PR/MR is merged into the protected default branch (`main`).

---

## 6. Changelog Updated

- `CHANGELOG.md` is updated if the change is user-visible.
- The update follows [Keep a Changelog](https://keepachangelog.com/) format.

---

## 7. Documentation Updated (if applicable)

- If the change affects setup, deployment, or API usage, the root `README.md` or relevant documentation is updated.

---

## 8. Quality Requirements and QRTs

- Relevant quality requirements (`QR-001`, `QR-002`, `QR-003`) are satisfied.
- Relevant automated quality requirement tests (`QRT-001`, `QRT-002`, `QRT-003`) pass.
- If a quality requirement is not applicable, it must be explicitly documented with rationale.

---

## 9. Work Status

- The issue Work Status is updated to `Done` in the issue tracker.

---

## 10. Definition of Done Applies to All PBIs

This Definition of Done applies to **all PBIs** — user stories, technical tasks, infrastructure work, testing tasks, and documentation work.

---

### ⚠️ Important

A PBI may be marked `Done` **only** when **all** of the above conditions are satisfied. If a condition cannot be met (e.g., a test is not feasible), the team must document the exception and obtain TA approval.
