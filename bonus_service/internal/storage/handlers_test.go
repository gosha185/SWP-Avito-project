package storage_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bonus-service/internal/handlers"
	"bonus-service/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := newTestDB(t)
	svc := service.NewBonusService(db)
	handler := handlers.NewAPIHandler(svc)
	router := handlers.NewRouter(handler)
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return srv
}

func post(t *testing.T, srv *httptest.Server, path string, body any, key string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, srv.URL+path, bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func get(t *testing.T, srv *httptest.Server, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(srv.URL + path)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(dst))
}

// --- Accrual ---

func TestHandler_Accrual_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": userID,
		"amount":  500,
		"days":    30,
	}, uuid.NewString())

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.NotZero(t, body["batch_id"])
}

func TestHandler_Accrual_MissingIdempotencyKey(t *testing.T) {
	srv := newTestServer(t)

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": uuid.New(),
		"amount":  100,
		"days":    7,
	}, "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func TestHandler_Accrual_InvalidBody(t *testing.T) {
	srv := newTestServer(t)

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"amount": -1,
		"days":   0,
	}, uuid.NewString())

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	resp.Body.Close()
}

func TestHandler_Accrual_DuplicateKey(t *testing.T) {
	srv := newTestServer(t)
	key := uuid.NewString()
	body := map[string]any{"user_id": uuid.New(), "amount": 100, "days": 7}

	resp1 := post(t, srv, "/v1/accrual", body, key)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	resp1.Body.Close()

	resp2 := post(t, srv, "/v1/accrual", body, key)
	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
	resp2.Body.Close()
}

// --- Hold ---

func TestHandler_Hold_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": userID, "amount": 1000, "days": 30,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = post(t, srv, "/v1/hold", map[string]any{
		"user_id":  userID,
		"order_id": uuid.New(),
		"amount":   300,
		"hours":    24,
	}, uuid.NewString())

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.NotZero(t, body["hold_id"])
}

func TestHandler_Hold_InsufficientBalance(t *testing.T) {
	srv := newTestServer(t)

	resp := post(t, srv, "/v1/hold", map[string]any{
		"user_id":  uuid.New(),
		"order_id": uuid.New(),
		"amount":   500,
		"hours":    24,
	}, uuid.NewString())

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	resp.Body.Close()
}

func TestHandler_Hold_MissingIdempotencyKey(t *testing.T) {
	srv := newTestServer(t)

	resp := post(t, srv, "/v1/hold", map[string]any{
		"user_id": uuid.New(), "order_id": uuid.New(), "amount": 100, "hours": 1,
	}, "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

// --- Confirm ---

func TestHandler_ConfirmHold_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()
	orderID := uuid.New()

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": userID, "amount": 1000, "days": 30,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = post(t, srv, "/v1/hold", map[string]any{
		"user_id": userID, "order_id": orderID, "amount": 200, "hours": 24,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	path := fmt.Sprintf("/v1/users/%s/holds/%s/confirm", userID, orderID)
	resp = post(t, srv, path, nil, uuid.NewString())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.NotZero(t, body["hold_id"])
}

func TestHandler_ConfirmHold_MissingKey(t *testing.T) {
	srv := newTestServer(t)
	path := fmt.Sprintf("/v1/users/%s/holds/%s/confirm", uuid.New(), uuid.New())
	resp := post(t, srv, path, nil, "")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

// --- Cancel ---

func TestHandler_CancelHold_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()
	orderID := uuid.New()

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": userID, "amount": 1000, "days": 30,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = post(t, srv, "/v1/hold", map[string]any{
		"user_id": userID, "order_id": orderID, "amount": 200, "hours": 24,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	path := fmt.Sprintf("/v1/users/%s/holds/%s/cancel", userID, orderID)
	resp = post(t, srv, path, nil, uuid.NewString())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.NotZero(t, body["hold_id"])
}

// --- Balance ---

func TestHandler_GetBalance_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()

	resp := post(t, srv, "/v1/accrual", map[string]any{
		"user_id": userID, "amount": 750, "days": 5,
	}, uuid.NewString())
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	path := fmt.Sprintf("/v1/users/%s/balance?days=10", userID)
	resp = get(t, srv, path)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.Equal(t, float64(750), body["available"])
	assert.Equal(t, float64(750), body["expiring"])
}

func TestHandler_GetBalance_NewUser(t *testing.T) {
	srv := newTestServer(t)
	path := fmt.Sprintf("/v1/users/%s/balance?days=30", uuid.New())
	resp := get(t, srv, path)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	assert.Equal(t, float64(0), body["available"])
	resp.Body.Close()
}

func TestHandler_GetBalance_MissingDays(t *testing.T) {
	srv := newTestServer(t)
	path := fmt.Sprintf("/v1/users/%s/balance", uuid.New())
	resp := get(t, srv, path)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	resp.Body.Close()
}

// --- History ---

func TestHandler_GetHistory_Success(t *testing.T) {
	srv := newTestServer(t)
	userID := uuid.New()

	for i := 0; i < 3; i++ {
		resp := post(t, srv, "/v1/accrual", map[string]any{
			"user_id": userID, "amount": 100, "days": 30,
		}, uuid.NewString())
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	path := fmt.Sprintf("/v1/users/%s/history?limit=10&offset=0", userID)
	resp := get(t, srv, path)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	txs := body["transactions"].([]any)
	assert.Len(t, txs, 3)
}

func TestHandler_GetHistory_MissingParams(t *testing.T) {
	srv := newTestServer(t)
	path := fmt.Sprintf("/v1/users/%s/history", uuid.New())
	resp := get(t, srv, path)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	resp.Body.Close()
}
