Learning points

1. Formulating verifiable requirements. The initial requirement “bonuses not lost during failures” is impossible to test absolutely. The customer helped us refine it into “no duplicate accruals on retries using idempotency keys”, which is concrete and testable.

2. Prioritisation with MoSCoW. We learned that a feature can be technically important but still be a “Should” or “Could” from the business perspective. For example, showing points that will expire soon seemed useful, but the customer considered it non‑essential for MVP v1, while showing the current balance is a clear Must Have.

3. User roles in an API product. Not all user stories are about end‑users. The main consumer of our API is another service.

Validated assumptions

- Users need to see their available balance before spending points. Customer said it’s "Must Have" 
- Showing points expiring soon is critical for MVP. Rejected: customer moved it to Could Have
- Idempotency keys are the right solution for duplicate prevention. Confirmed: customer recommended this approach
- The system must support partial hold release. Rejected: customer said only full release is needed
- Spending order should be FEFO. Confirmed: customer suggested it and team agreed

Needs clarification

No major unresolved questions after the customer meeting. All technical details (idempotency, TTL, FEFO) are agreed and will be implemented.

Planned response

- MVP v1 scope now contains: US‑06 (view held points), US‑07 (view available balance).
- Technical implementation:
    Implement idempotency keys for all state‑changing endpoints (accrual, hold, confirm, release).
    Add a background job to automatically release expired holds after a configurable TTL.
    Implement FEFO spending logic.
