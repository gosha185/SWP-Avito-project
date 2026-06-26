# User Acceptance Tests

This document defines end-user-facing acceptance test scenarios for the Avito Bonus Points Service.

---

## UAT-001: View current available balance

**Status:** Active

**User goal:** As an end user, I want to see my current available bonus balance so that I know how many points I can spend.

**Preconditions:**
- User is registered in the system
- User has some bonus points (e.g., 150 points available)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/balance` request.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK response.
- Response body contains `balance` field with the correct number of points.
- Response time < 500ms.

**Execution result:** (to be filled after UAT session)

**Customer comments:** (to be filled after UAT session)

---

## UAT-002: View points held for a concrete order

**Status:** Active

**User goal:** As an end user, I want to see how many points are held for a concrete order so that I can analyse the number of points awarded.

**Preconditions:**
- User is registered in the system
- User has an active order with held points (e.g., 50 points held)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/holds/{order_id}` request.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK response.
- Response body contains `held` field with the correct number of points.
- Response time < 500ms.

**Execution result:** (to be filled after UAT session)

**Customer comments:** (to be filled after UAT session)

---

## UAT-003: Accrue points with TTL

**Status:** Active

**User goal:** As an end user, I want to receive bonus points after completing a qualifying action so that I am rewarded for my activity.

**Preconditions:**
- User is registered in the system
- Qualifying action trigger is available (e.g., purchase completed)

**Step-by-step instructions:**
1. Trigger a qualifying action (e.g., complete a test purchase).
2. Wait for the accrual process to complete.
3. Send `GET /v1/users/{user_id}/balance` request.
4. Observe the response.

**Expected outcome:**
- HTTP 200 OK response.
- User's balance increased by the expected amount.
- Points have a valid TTL (expiration timestamp).
- Response time < 500ms.

**Execution result:** (to be filled after UAT session)

**Customer comments:** (to be filled after UAT session)
