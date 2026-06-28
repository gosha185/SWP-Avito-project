# User Acceptance Tests

This document defines end-user-facing acceptance test scenarios for the Avito Bonus Points Service.

---

## UAT-001: View current available balance

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-001 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to see my current available bonus balance so that I know how many points I can spend. |

**Preconditions:**
- User is registered in the system
- User has bonus points (e.g., 500 points)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/balance`.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK
- `balance` field with correct number of points

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "That looks correct. We placed 200 points on hold, so there should be 300 available points remaining. That matches the expected behavior."

**Resulting PBIs or issues:** None

---

## UAT-002: View points held for a concrete order

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-002 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to see how many points are held for a concrete order so that I can analyse the number of points awarded. |

**Preconditions:**
- User is registered in the system
- User has an active hold (e.g., 200 points held for an order)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/holds`.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK
- List of active holds with correct amounts

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "If we open the list of active holds, we can see that all 700 points have been successfully reserved."

**Resulting PBIs or issues:** None

---

## UAT-003: Award points with TTL

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-003 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to receive bonus points after completing a qualifying action. |

**Preconditions:**
- User is registered in the system

**Step-by-step instructions:**
1. Trigger a qualifying action (e.g., purchase).
2. Send `POST /v1/accrue` with `user_id`, `amount`, `ttl_seconds`.
3. Send `GET /v1/users/{user_id}/balance`.
4. Observe the response.

**Expected outcome:**
- HTTP 200 OK
- Balance increased by the expected amount
- Points have a valid expiration timestamp

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "Works correctly."

**Resulting PBIs or issues:** None

---

## UAT-004: Place points on hold

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-004 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to place bonus points on hold for a specific order. |

**Preconditions:**
- User is registered in the system
- User has sufficient available points (e.g., 500 points)

**Step-by-step instructions:**
1. Send `POST /v1/hold` with `user_id`, `order_id`, `amount`.
2. Send `GET /v1/users/{user_id}/holds`.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK
- Hold appears in the holds list
- Available balance decreased by the held amount

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "All 700 points have been successfully reserved."

**Resulting PBIs or issues:** None

---

## UAT-005: Confirm a hold

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-005 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to confirm a hold and complete the points deduction. |

**Preconditions:**
- User has an active hold (e.g., 200 points held)

**Step-by-step instructions:**
1. Send `POST /v1/confirm` with `user_id`, `order_id`.
2. Send `GET /v1/users/{user_id}/holds`.
3. Send `GET /v1/users/{user_id}/balance`.

**Expected outcome:**
- HTTP 200 OK
- Hold removed from holds list
- Balance remains unchanged (already deducted at hold creation)

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "That's exactly how it should work."

**Resulting PBIs or issues:** None

---

## UAT-006: Cancel a hold

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-006 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to cancel a hold and release the points back to my balance. |

**Preconditions:**
- User has an active hold (e.g., 200 points held)

**Step-by-step instructions:**
1. Send `POST /v1/cancel` with `user_id`, `order_id`.
2. Send `GET /v1/users/{user_id}/holds`.
3. Send `GET /v1/users/{user_id}/balance`.

**Expected outcome:**
- HTTP 200 OK
- Hold removed from holds list
- Available balance increased by the cancelled amount

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "We held 200 points, cancelled the hold, and those points became available again."

**Resulting PBIs or issues:** None

---

## UAT-007: Cancel an already cancelled hold (error handling)

| Field | Value |
|-------|-------|
| **Scenario ID** | UAT-007 |
| **Scenario status** | Active |
| **User goal** | As an end user, I want to receive a clear error when trying to cancel a hold that no longer exists. |

**Preconditions:**
- A hold was previously cancelled

**Step-by-step instructions:**
1. Send `POST /v1/cancel` with `user_id`, `order_id` for an already cancelled hold.
2. Observe the response.

**Expected outcome:**
- HTTP 400 Bad Request
- Clear error message: "Hold already cancelled"

**Execution result (25 June 2026):** ❌ Failed

**Customer comments:** "This definitely needs improvement. I'd expect HTTP 400 (Bad Request) instead. The request itself is invalid rather than the server encountering an internal error."

**Resulting PBIs or issues:** [#174](https://github.com/gosha185/SWP-Avito-project/issues/177)

---

## Summary of UAT Results (25 June 2026)

| UAT | Scenario | Status |
|-----|----------|--------|
| UAT-001 | View current available balance | ✅ Passed |
| UAT-002 | View points held for a concrete order | ✅ Passed |
| UAT-003 | Award points with TTL | ✅ Passed |
| UAT-004 | Place points on hold | ✅ Passed |
| UAT-005 | Confirm a hold | ✅ Passed |
| UAT-006 | Cancel a hold | ✅ Passed |
| UAT-007 | Cancel an already cancelled hold (error handling) | ❌ Failed — see [#174](https://github.com/gosha185/SWP-Avito-project/issues/174) |
