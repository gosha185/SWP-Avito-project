# HTTP Layer (Handlers) Documentation

The `internal/handlers` package implements the HTTP interface layer. It converts HTTP requests into service calls, handles validation, serializes responses, and manages error handling according to OpenAPI specification.

## Structure

The package consists of:
- **APIHandler** - main handler struct with all HTTP handlers
- **DTOs** - request/response data transfer objects (`dto.go`)
- **Error helpers** - HTTP error response helpers (`errors.go`)
- **JSON helpers** - JSON parsing/writing utilities (`helpers.go`)
- **Router** - HTTP route definitions (`router.go`)

## APIHandler

### NewAPIHandler

**Purpose:** Creates a new APIHandler with a service instance.

**Parameters:**
- `s *service.BonusService` - service layer instance

**Returns:**
- `*APIHandler` - new handler instance

**Usage:**
```go
handler := handlers.NewAPIHandler(serviceInstance)
router := handlers.NewRouter(handler)
```

## Handlers

Each handler follows the same pattern:
1. Extract and validate required headers (Idempotency-Key when needed)
2. Parse and validate request body/path parameters
3. Convert to service layer model (`LedgerEntry`)
4. Call appropriate service method
5. Map service errors to HTTP status codes
6. Return JSON response

### AccrualHandler

**Purpose:** Handle POST /accrual - create new points batch.

**HTTP Method:** `POST /accrual`

**Request Headers:**
- `Idempotency-Key: <uuid>` (required)

**Request Body:**
```json
{
  "user_id": "uuid",
  "amount": 500,
  "days": 30
}
```

**Validation:**
- `user_id` must be valid UUID
- `amount` must be positive
- `days` must be positive

**Service Call:** `service.Accrue()`

**Success Response:** `201 Created`
```json
{
  "batch_id": 1
}
```

**Error Responses:**
- `400 Bad Request` - missing Idempotency-Key or invalid request body
- `422 Unprocessable Entity` - validation failed
- `409 Conflict` - idempotency key already used
- `500 Internal Server Error` - service error

---

### CreateHoldHandler

**Purpose:** Handle POST /hold - reserve points for an order.

**HTTP Method:** `POST /hold`

**Request Headers:**
- `Idempotency-Key: <uuid>` (required)

**Request Body:**
```json
{
  "user_id": "uuid",
  "order_id": "uuid",
  "amount": 200,
  "hours": 24
}
```

**Validation:**
- `user_id` must be valid UUID
- `order_id` must be valid UUID
- `amount` must be positive
- `hours` must be positive

**Service Call:** `service.Hold()`

**Success Response:** `201 Created`
```json
{
  "hold_id": 7
}
```

**Error Responses:**
- `400 Bad Request` - missing Idempotency-Key or invalid request body
- `422 Unprocessable Entity` - validation failed
- `409 Conflict` - idempotency key already used, insufficient balance, or order already held
- `500 Internal Server Error` - service error

---

### ConfirmHoldHandler

**Purpose:** Handle POST /users/{user_id}/holds/{order_id}/confirm - finalize a reservation.

**HTTP Method:** `POST /users/{user_id}/holds/{order_id}/confirm`

**Request Headers:**
- `Idempotency-Key: <uuid>` (required)

**Path Parameters:**
- `user_id` - user UUID
- `order_id` - order UUID

**Validation:**
- Both path parameters must be valid UUIDs

**Service Call:** `service.ConfirmWithdraw()`

**Success Response:** `200 OK`
```json
{
  "hold_id": 7
}
```

**Error Responses:**
- `400 Bad Request` - missing Idempotency-Key or invalid UUIDs
- `422 Unprocessable Entity` - validation failed
- `409 Conflict` - idempotency key already used
- `500 Internal Server Error` - service error

---

### CancelHoldHandler

**Purpose:** Handle POST /users/{user_id}/holds/{order_id}/cancel - release a reservation.

**HTTP Method:** `POST /users/{user_id}/holds/{order_id}/cancel`

**Request Headers:**
- `Idempotency-Key: <uuid>` (required)

**Path Parameters:**
- `user_id` - user UUID
- `order_id` - order UUID

**Validation:**
- Both path parameters must be valid UUIDs

**Service Call:** `service.CancelHold()`

**Success Response:** `200 OK`
```json
{
  "hold_id": 7
}
```

**Error Responses:**
- `400 Bad Request` - missing Idempotency-Key or invalid UUIDs
- `422 Unprocessable Entity` - validation failed
- `409 Conflict` - idempotency key already used
- `500 Internal Server Error` - service error

---

### GetBalanceHandler

**Purpose:** Handle GET /users/{user_id}/balance - get available points.

**HTTP Method:** `GET /users/{user_id}/balance`

**Path Parameters:**
- `user_id` - user UUID

**Validation:**
- `user_id` must be valid UUID

**Service Call:** `service.GetAvailablePoints()`

**Success Response:** `200 OK`
```json
{
  "available": 500
}
```

**Error Responses:**
- `400 Bad Request` - invalid UUID
- `500 Internal Server Error` - service error

---

### GetExpirationsHandler

**Purpose:** Handle GET /users/{user_id}/balance/expirations - get points expiring within N days.

**HTTP Method:** `GET /users/{user_id}/balance/expirations`

**Path Parameters:**
- `user_id` - user UUID

**Query Parameters:**
- `days` - number of days to look ahead (required, positive integer)

**Validation:**
- `user_id` must be valid UUID
- `days` must be provided and parsable as integer

**Service Call:** `service.GetExpiringAvailablePoints()`

**Success Response:** `200 OK`
```json
{
  "expiring": 150
}
```

**Error Responses:**
- `400 Bad Request` - invalid UUID or missing/invalid days parameter
- `500 Internal Server Error` - service error

---

### GetHistoryHandler

**Purpose:** Handle GET /users/{user_id}/history - get paginated transaction history.

**HTTP Method:** `GET /users/{user_id}/history`

**Path Parameters:**
- `user_id` - user UUID

**Query Parameters:**
- `limit` - page size (required, positive integer)
- `offset` - number of items to skip (required, non-negative integer)

**Validation:**
- `user_id` must be valid UUID
- `limit` must be provided and parsable as positive integer
- `offset` must be provided and parsable as non-negative integer

**Service Call:** `service.GetHistory()`

**Success Response:** `200 OK`
```json
{
  "transactions": [
    {
      "operation_type": "accrual",
      "amount": 500,
      "created_at": "2023-01-01T12:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request` - invalid UUID or missing/invalid pagination parameters
- `500 Internal Server Error` - service error

---

### GetHoldsHandler

**Purpose:** Handle GET /users/{user_id}/holds - get total held points.

**HTTP Method:** `GET /users/{user_id}/holds`

**Path Parameters:**
- `user_id` - user UUID

**Validation:**
- `user_id` must be valid UUID

**Service Call:** `service.GetHeld()`

**Success Response:** `200 OK`
```json
{
  "amount": 150
}
```

**Error Responses:**
- `400 Bad Request` - invalid UUID
- `500 Internal Server Error` - service error

## DTOs (Data Transfer Objects)

All request/response structures are defined in `dto.go`:

- `accrualRequest` - POST /accrual request body
- `accrualResponse` - POST /accrual success response
- `createHoldRequest` - POST /hold request body
- `createHoldResponse` - POST /hold success response
- `confirmHoldResponse` - confirm hold success response
- `cancelHoldResponse` - cancel hold success response
- `getBalanceResponse` - GET /balance success response
- `getExpirationsResponse` - GET /expirations success response
- `transaction` - single history entry
- `GetHistoryResponse` - GET /history success response
- `GetHoldResponse` - GET /holds success response

## Error Handling

Error helpers map different error types to appropriate HTTP responses:

- `errorResponse()` - generic error response
- `serverErrorResponse()` - 500 Internal Server Error
- `notFoundResponse()` - 404 Not Found
- `methodNotAllowedResponse()` - 405 Method Not Allowed
- `badRequestResponse()` - 400 Bad Request
- `failedValidationResponse()` - 422 Unprocessable Entity

**Service Error Mapping:**
- `storage.ErrLedgerDuplicate` → `409 Conflict`
- `storage.ErrInsufficientBalance` → `409 Conflict`
- `storage.ErrOrderAlreadyHeld` → `409 Conflict`
- `storage.ErrHoldNotFound` → `400 Bad Request`
- All other service errors → `500 Internal Server Error`

## JSON Helpers

### readJSON()

**Purpose:** Parse JSON request body with validation.

**Features:**
- Limits body size to 10KB
- Disallows unknown fields
- Returns descriptive error messages for:
  - Malformed JSON
  - Wrong JSON type
  - Empty body
  - Unknown fields
  - Body too large

### writeJSON()

**Purpose:** Write JSON response with proper headers.

**Features:**
- Sets `Content-Type: application/json`
- Appends newline to output
- Supports custom headers

### readIDParam()

**Purpose:** Parse UUID from path parameter "id".

**Note:** Currently not used - handlers parse UUIDs directly from path values.

## Router

### NewRouter()

**Purpose:** Create HTTP router with all route definitions.

**Routes:**
```
POST   /accrual
POST   /hold
POST   /users/{user_id}/holds/{order_id}/confirm
POST   /users/{user_id}/holds/{order_id}/cancel
GET    /users/{user_id}/balance
GET    /users/{user_id}/balance/expirations
GET    /users/{user_id}/history
GET    /users/{user_id}/holds
```

**Returns:** `*http.ServeMux` configured with all handlers

## Conventions

1. **Idempotency-Key Header:** Required for all write operations (POST endpoints)
2. **UUID Validation:** All UUIDs are validated before service calls
3. **JSON Serialization:** All responses are JSON with proper Content-Type
4. **Error Responses:** Consistent error format: `{"error": "message"}`
5. **Validation Errors:** Return field-level errors for validation failures
6. **Path Parameters:** Extracted using `r.PathValue()` (Go 1.22+)
7. **Service Layer Mapping:** Handlers convert between HTTP and service models

## Flow Example (POST /hold)

```
1. Client → HTTP Request with Idempotency-Key header and JSON body
2. Handler → Validate Idempotency-Key presence
3. Handler → Parse JSON into createHoldRequest DTO
4. Handler → Validate fields using validator
5. Handler → Create LedgerEntry from request
6. Handler → Call service.Hold() with LedgerEntry
7. Service → Business logic and storage operations
8. Handler → Map service errors to HTTP status codes
9. Handler → Create createHoldResponse DTO
10. Handler → Write JSON response
11. Client ← HTTP Response with hold_id or error
```