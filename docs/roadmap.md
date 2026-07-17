# Product Roadmap

## Sprint 1 – MVP v1 (Completed)

**Dates:** 15–21 June 2026  
**Sprint Goal:** Deliver core bonus operations: users can view their available balance and see points held for orders. The system supports idempotent accrual, two-phase redemption with FEFO, and automatic release of stale holds.

**Milestone:** [Sprint 1 – MVP v1](https://github.com/gosha185/SWP-Avito-project/milestone/1)

**Delivered:**
- US-06: View points held for a concrete order
- US-07: View current available point balance
- 18 technical and infrastructure PBIs completed

---

## Sprint 2 – Quality and Testing (Completed)

**Dates:** 22–28 June 2026  
**Sprint Goal:** Stabilise the product and establish quality foundations: complete the remaining MVP v1 features, implement automated testing and CI, and ensure the application runs reliably and passes all tests.

**Milestone:** [Sprint 2 – Quality and Testing](https://github.com/gosha185/SWP-Avito-project/milestone/2)

**Completed:**
- Quality requirements (QR-001 to QR-005)
- Quality requirement tests (QRT-001 to QRT-005)
- Unit and integration tests (coverage 53.2%)
- CI pipeline with linting, tests, coverage, QRTs, additional QA
- UAT scenarios (7 scenarios, 6 passed, 1 failed — fixed in Sprint 3)
- Updated Definition of Done
- Week 4 reports and documentation

---

## Sprint 3 – MVP v2 (Completed)

**Dates:** 29 June – 5 July 2026  
**Sprint Goal:** Deliver MVP v2: complete remaining features, document architecture (static, dynamic, deployment views), create 3 ADRs, and document development process with git workflow diagram.

**Milestone:** [Sprint 3 – MVP v2](https://github.com/gosha185/SWP-Avito-project/milestone/3)

**Completed:**
- TTL workers for expired points and holds (#9, #32)
- View points expiring in given time (#8)
- View all current held points (#23)
- HTTP 500 → HTTP 400 error handling (#177)
- Architecture documentation (static, dynamic, deployment views)
- 3 ADRs created (ADR-001, ADR-002, ADR-003)
- Development process documentation with git workflow diagram
- UAT scenarios (9 scenarios, all passed)
- Sprint Review with customer

---

## Sprint 4 – Trial Release and Handover Preparation (Current)

**Dates:** 6–12 July 2026  
**Sprint Goal:** Deliver a stable trial release, complete customer-facing documentation (README, customer-handover, CONTRIBUTING, AGENTS), and conduct a transition-readiness meeting with the customer.

**Milestone:** [Sprint 4 – Trial Release](https://github.com/gosha185/SWP-Avito-project/milestone/4)

**Planned items:**

| # | Title | Type | SP | Status |
|---|-------|------|----|--------|
| #261 | Docs: Update roadmap for Sprint 4 | Documentation | 1 | Done |
| #262 | Docs: Update roadmap for Sprint 5 | Documentation | 1 | Done |
| #263 | Docs: Update README.md for customer handover | Documentation | 3 | Done |
| #264 | Docs: Create customer-handover.md | Documentation | 5 | Done |
| #265 | Docs: Create CONTRIBUTING.md | Documentation | 2 | Done |
| #266 | Docs: Create AGENTS.md | Documentation | 2 | Done |
| #267 | Infra: Deploy Week 6 trial release | Infrastructure | 3 | Done |
| #268 | Docs: Update CHANGELOG.md for Week 6 | Documentation | 1 | Done |
| #269 | Docs: Add Week 6 report index | Documentation | 5 | Done |
| #270 | Docs: Add Week 6 reflection | Documentation | 2 | Done |
| #271 | Docs: Add Week 6 retrospective | Documentation | 2 | Done |
| #272 | Docs: Add LLM usage report for Week 6 | Documentation | 1 | Done |
| #273 | Docs: Add Week 6 sprint-review-summary and meeting transcript | Documentation | 3 | Done |

**Focus areas:**
- Deploy trial release for customer testing
- Update README.md as the main entry point
- Create customer-handover.md with handover status
- Create CONTRIBUTING.md and AGENTS.md
- Conduct transition-readiness meeting with customer
- Gather customer feedback for Week 7 follow-up

---

## Sprint 5 – MVP v3 (Completed)

**Dates:** 13–19 July 2026  
**Sprint Goal:** Address trial feedback, complete final transition, deliver MVP v3, and prepare for Demo Day.

**Milestone:** [Sprint 5 – MVP v3](https://github.com/gosha185/SWP-Avito-project/milestone/5)

**Completed:**
- Final handover meeting with customer
- Handover status confirmed: `Ready for independent use`
- Customer confirmation: `Accepted`
- Bug fixes (error codes, rare edge cases, ledger data insertion)
- Code fully commented and documented
- User stories marked as implemented
- Week 7 reports: README, sprint-review-transcript, sprint-review-summary, reflection, retrospective, llm-report
- Public sanitized demo video for MVP v3
- SemVer release for MVP v3 created
- CHANGELOG.md updated
- Demo Day presentation prepared
