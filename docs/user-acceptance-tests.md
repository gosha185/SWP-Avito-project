# User Acceptance Tests

This document defines end-user-facing acceptance test scenarios for the Avito Bonus Points Service.

---

## UAT-001: View current available balance

**Scenario ID:** UAT-001  
**Scenario status:** Active  
**User goal:** As an end user, I want to see my current available bonus balance so that I know how many points I can spend.

**Preconditions:**
- User is registered in the system
- User has some bonus points (e.g., 500 points available)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/balance` request.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK response.
- Response body contains `balance` field with the correct number of points.

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "That looks correct. We placed 200 points on hold, so there should be 300 available points remaining. That matches the expected behavior."

**Resulting PBIs or issues:** None.

---

## UAT-002: View points held for a concrete order

**Scenario ID:** UAT-002  
**Scenario status:** Active  
**User goal:** As an end user, I want to see how many points are held for a concrete order so that I can analyse the number of points awarded.

**Preconditions:**
- User is registered in the system
- User has an active order with held points (e.g., 200 points held)

**Step-by-step instructions:**
1. Authenticate as a test user.
2. Send `GET /v1/users/{user_id}/holds` request.
3. Observe the response.

**Expected outcome:**
- HTTP 200 OK response.
- Response body contains list of active holds.

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "If we open the list of active holds, we can see that all 700 points have been successfully reserved."

**Resulting PBIs or issues:** None.

---

## UAT-003: Accrue points with TTL and hold operations

**Scenario ID:** UAT-003  
**Scenario status:** Active  
**User goal:** As an end user, I want to receive bonus points, have them placed on hold, and then confirm or cancel the hold.

**Preconditions:**
- User is registered in the system
- Qualifying action trigger is available

**Step-by-step instructions:**
1. Trigger a qualifying action (e.g., complete a test purchase).
2. Send `POST /v1/accrue` to award points.
3. Send `POST /v1/hold` to place points on hold.
4. Send `POST /v1/confirm` or `POST /v1/cancel` to finalise.
5. Verify balance and hold list.

**Expected outcome:**
- Points awarded correctly.
- Hold created and visible.
- Confirm/cancel works as expected.
- Balance updated correctly.

**Execution result (25 June 2026):** ✅ Passed

**Customer comments:** "Everything works correctly in your implementation. I don't see any issues here."

**Resulting PBIs or issues:** None.

---

## Additional Findings from UAT Session

| Issue | Severity | Action | Resulting PBI |
|-------|----------|--------|---------------|
| Cancelling a hold that was already cancelled returns HTTP 500 instead of HTTP 400 | Minor | Fix error handling to return 400 Bad Request | [#XXX](https://github.com/gosha185/SWP-Avito-project/issues/XXX) |
| FEFO (first-expiring first-out) logic implemented | N/A | Verified — points with earliest expiration are spent first | None |
