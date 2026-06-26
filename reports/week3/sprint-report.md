# Sprint Report – Week 3

## Sprint Goal

Deliver core bonus operations: users can view their available balance and see points held for orders. The system supports idempotent accrual, two-phase redemption with FEFO, and automatic release of stale holds.

---

## Sprint Backlog

| # | Title | Assignee | Status | Story Points |
|---|-------|----------|--------|--------------|
| #4 | US-07: View current available point balance | Georgii, Ekaterina, Stepan | Done | 1 |
| #5 | US-06: View points held for a concrete order | Georgii, Ekaterina, Stepan | Done | 1 |
| #8 | Technical: Make possible to see points expiring in giving time | Georgii, Ekaterina, Stepan | Next Sprint | 13 |
| #9 | Technical: Make the TTL worker manage expired points batches | Ekaterina | Next Sprint | 8 |
| #10 | Concurrency tests | Ivan | Done | 8 |
| #11 | Technical: Configure PostgreSQL storage | Ekaterina | Done | 5 |
| #12 | Public deployment of MVP v1 | Stepan | Next Sprint | 2 |
| #13 | Update OpenAPI / Swagger UI | Stepan | Done | 8 |
| #14 | Update Postman collection | Stepan | Next Sprint | 5 |
| #15 | SemVer release v1.0.0 | Stepan | Next Sprint | 2 |
| #16 | CHANGELOG.md update | Stepan | Next Sprint | 1 |
| #23 | Technical: Make possible to see all current held points | Georgii, Ekaterina, Stepan | Next Sprint | 1 |
| #24 | Technical: Write an algorithm to accrue points | Georgii | Done | 3 |
| #25 | Technical: Write an algorithm to hold points | Georgii | Done | 3 |
| #26 | Technical: Write an algorithm to cancel points holding by the order id | Georgii | Done | 3 |
| #27 | Technical: Write an algorithm to withdraw points from hold | Georgii | Done | 3 |
| #28 | Technical: Write http layer for processing given commands | Stepan | Done | 5 |
| #29 | Technical: Write validation and idempotency check | Stepan | Done | 5 |
| #30 | Technical: Make every mutating operation written down in ledger | Stepan | Done | 5 |
| #31 | Technical: Write error handling for errors which may be solved on our side | Stepan | Done | 5 |
| #32 | Technical: Make the TTL worker manage expired held points | Ekaterina | Next Sprint | 8 |
| #33 | Technical: Write interactions with balances at the database level | Ekaterina | Done | 3 |
| #34 | Technical: Write interactions with batches of points at the database level | Ekaterina | Done | 3 |
| #35 | Technical: Write interactions with held points at the database level | Ekaterina | Done | 3 |
| #36 | Technical: Write interactions with ledger at the database level | Ekaterina | Done | 3 |

**Total Sprint Size:** 101 Story Points  
**Completed:** 18 tasks  
**Moved to Next Sprint:** 7 tasks (#8, #9, #12, #14, #15, #16, #23, #32)

---

## Product Backlog Size

**Total Product Backlog Size:** 114 Story Points (27 tasks)

---

## Changes Since Assignment 2

### User Stories
- **US-07** (View current available point balance) – promoted to Must Have, added to MVP v1
- **US-06** (View points held for a concrete order) – remains Must Have
- **US-04** (Warning about exceeding point limit) – moved to Should Have (not in MVP v1)
- **US-08** (Admin view transaction history) – moved to Should Have (not in MVP v1)

### Technical Tasks
- Large tasks were split into smaller, estimable items:
  - "Accrual with TTL + idempotency" → split into #24, #28, #29, #30, #31
  - "Two-phase redemption + FEFO" → split into #25, #26, #27
- New tasks added: #8, #9, #23, #32, #33, #34, #35, #36

### Customer Feedback Addressed
- Showing only available balance is **Must Have** – US-07 added
- Showing points expiring soon moved to **Could Have** (implemented as #8)
- Idempotency requirement refined to "no duplicate accruals" – captured in #29
- FEFO (first-expiring first-out) spending order – captured in #26, #27

---

## Links
- [Historical User Stories (Assignment 2)](../week2/user-stories.md)
- [Current User Stories Index](../../docs/user-stories.md)
- [Product Backlog Board](https://github.com/users/gosha185/projects/1)
- [Sprint Backlog Board](https://github.com/users/gosha185/projects/2)
- [Sprint Milestone](https://github.com/gosha185/SWP-Avito-project/milestones/1)
- [MVP v1 Scope](https://github.com/gosha185/SWP-Avito-project/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Sprint+1+–+MVP+v1%22)

---

## MVP v1 Scope

The following 25 tasks were selected for MVP v1 (all Must Have).

**Delivered (18 tasks):**
- US-07, US-06, #10, #11, #13, #24, #25, #26, #27, #28, #29, #30, #31, #33, #34, #35, #36

**Moved to Next Sprint (7 tasks):**
- #8, #9, #12, #14, #15, #16, #23, #32

---

## PBI Types and Statuses

| Type | Count | Statuses |
|------|-------|----------|
| User Story | 4 | Active |
| Technical PBI | 18 | Active |
| Infrastructure PBI | 5 | Active |

**Work Status:** To Do → In Progress → Review → Done

---

## Verification Evidence

Each completed MVP v1 PBI was verified through:
- Acceptance criteria checked in PR description
- Code review by another team member
- PR merged into protected default branch
- All CI checks passing
- CHANGELOG.md updated for user-visible changes

---

## Current Product Status

**MVP v1:** 18 of 25 tasks completed (72%). 7 tasks moved to Sprint 2.

---

## Next Steps

1. Complete remaining MVP v1 tasks in Sprint 2: #8, #9, #12, #14, #15, #16, #23, #32
2. Split large tasks (#8, #9) into smaller subtasks (≤ 5 SP)
3. Pass all acceptance criteria
4. Create SemVer release v1.0.0
5. Prepare Assignment 3 submission

---

## Contribution Traceability

| Team Member | Issues Assigned | PRs Created | Reviews Done |
|-------------|-----------------|-------------|--------------|
| Georgii (gosha185) | #24, #25, #26, #27 | yes | yes |
| Ekaterina (deadnothingness) | #9, #11, #32, #33, #34, #35, #36 | yes | yes |
| Stepan (Stepan4ick) | #12, #13, #14, #15, #16, #28, #29, #30, #31 | yes | yes |
| Ivan (Laplace-mt) | #10 | yes | yes |
| Leilia (Leilia34) | #4, #5, #8, #23, documentation | yes | yes |
