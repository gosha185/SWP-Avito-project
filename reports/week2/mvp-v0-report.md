# MVP v0 — ITMO Ledger

## Purpose and description

ITMO Ledger is the runnable technical foundation for the bonus / points system.
It is a small HTTP API (Go + PostgreSQL) that stores one balance per user and
changes it through deposit and withdrawal transactions. The service is deployed
and runnable, exposes a health endpoint, and ships with a self-hosted Swagger UI
that renders the OpenAPI specification.

MVP v0 is a product foundation, not a finished user story: it establishes a
working, deployable API skeleton that later user stories build on.

## Access

| What | Link |
| --- | --- |
| Deployed Swagger UI (rendered OpenAPI) | http://10.93.26.189:8080/ |
| Deployed API base path | http://10.93.26.189:8080/v1 |
| Health endpoint | http://10.93.26.189:8080/v1/healthcheck |
| OpenAPI specification (repo) | ../../api/openapi.yaml |
| Postman collection (repo) | ../../api/postman_collection.json |
| Public Postman workspace (view-only) | https://www.postman.com/stepagrek07-2276551/itmo-ledger-api/overview |

**Access note.** The VM uses a private university IP address (`10.93.26.189`),
reachable from the university network.
<!-- TODO: confirm with the TA that university-network access is acceptable for
grading. If internet access is required, put a public tunnel (e.g. Cloudflare
Tunnel) in front of port 8080 and replace the links above with the public URL. -->

## Public video demonstration

<!-- TODO: add a public, sanitized video demonstration shorter than 2 minutes -->
- Video: TODO

## Relationship to the prototype and MVP v1 stories

The API artifacts (OpenAPI spec, Swagger UI, Postman collection) are the
interface prototype for the API product. MVP v0 implements the foundation behind
the following user stories:

<!-- TODO: reference the stable user-story IDs from reports/week2/user-stories.md
that this API foundation supports (for example: crediting points to a user, and
viewing a user's balance). -->
- US-XX: TODO
- US-YY: TODO

## Current limitations, placeholders, and mocks

- No authentication or authorization; any caller can change any balance.
- A withdrawal for a user that does not yet exist currently **creates** a balance
  equal to the withdrawal amount instead of rejecting it (known issue).
- On a database error during create/update, the HTTP response may be written
  twice (known issue, cosmetic).
- Only the two business endpoints plus the health endpoint are implemented; no
  listing, pagination, or transaction history.
- Amounts are stored as integers (no currency / decimal handling).

## Local setup

See the **Run with Docker** section in the root [`README.md`](../../README.md).
Short version:

    cd itmo_ledger
    cp .env.example .env          # then set POSTGRES_PASSWORD in .env
    docker compose up --build -d

After startup: Swagger UI at http://localhost:8080/, API at
http://localhost:8080/v1.

## Smoke-check scenario (repeatable)

**Goal:** confirm that MVP v0 is accessible and usable for its purpose — the
service is up, the health endpoint responds, Swagger UI loads, and a real
data-flow call works.

**Access:** use the deployed host `http://10.93.26.189:8080` (university
network), or `http://localhost:8080` after the local setup above.

**Steps and expected results:**

1. **Health endpoint.**

       curl http://10.93.26.189:8080/v1/healthcheck

   Expected: HTTP 200, body `{"status":"available"}`.

2. **Swagger UI.** Open `http://10.93.26.189:8080/` in a browser.
   Expected: the Swagger UI page renders and lists the endpoints
   (`/healthcheck`, `/transactions`, `/users/{id}/balance`).

3. **Create a balance (data-flow).**

       curl -X POST http://10.93.26.189:8080/v1/transactions \
         -d '{"user_id":"653F535D-10BA-4186-A05B-74493354F13B","amount":100,"type":"deposit"}'

   Expected: HTTP 201 (or 200 if the user already exists) and a balance object
   with `"amount": 100` (or the updated amount).

4. **Read the balance back.**

       curl http://10.93.26.189:8080/v1/users/653F535D-10BA-4186-A05B-74493354F13B/balance

   Expected: HTTP 200 and a balance object reflecting the amount from step 3.

A successful run of all four steps confirms MVP v0 is deployed and usable.
