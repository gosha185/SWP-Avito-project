package handlers

import (
	"bonus-service/internal/models"
	"bonus-service/internal/service"
	"bonus-service/internal/storage"
	"bonus-service/internal/validator"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// APIHandler implements all HTTP endpoints for the bonus service.
// It bridges HTTP requests to the service layer and handles JSON serialization.
type APIHandler struct {
	service *service.BonusService
}

// NewAPIHandler creates a new APIHandler with a service instance.
func NewAPIHandler(s *service.BonusService) *APIHandler {
	return &APIHandler{service: s}
}

//----------------------------- HANDLERS ------------------------------------------

// AccrualHandler handles POST /accrual requests to accrue bonus points for a user.
// Validates Idempotency-Key header, request body, and calls service.Accrue().
// Returns 201 with accrual batch ID on success.
// Error responses: 400 (bad request), 422 (validation), 409 (conflict - duplicate key), 500 (server error).
func (h *APIHandler) AccrualHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		h.badRequestResponse(w, r, errors.New("Idempotency-Key header is required"))
		return
	}

	var request accrualRequest
	if err := h.readJSON(w, r, &request); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(request.UserID != uuid.Nil, "user_id", "must be provided")
	v.Check(request.Amount > 0, "amount", "must be positive")
	v.Check(request.Days > 0, "days", "must be positive")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	entry := &models.LedgerEntry{
		UserID:        request.UserID,
		OperationType: models.OpAccrual,
		Amount:        request.Amount,
		ExternalKey:   key,
	}

	batchID, err := h.service.Accrue(r.Context(), entry, int64(request.Days))
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLedgerDuplicate):
			h.errorResponse(w, r, http.StatusConflict, "idempotency key already used for another operation")
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	response := accrualResponse{
		BatchId: batchID,
	}

	if err := h.writeJSON(w, http.StatusCreated, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// CreateHoldHandler handles POST /hold requests to reserve points for an order.
// Validates Idempotency-Key header, request body, and calls service.Hold().
// Returns 201 with hold ID on success.
// Error responses: 400 (bad request), 422 (validation), 409 (conflict - duplicate key, insufficient balance, order already held), 500 (server error).
func (h *APIHandler) CreateHoldHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		h.badRequestResponse(w, r, errors.New("Idempotency-Key header is required"))
		return
	}

	var request createHoldRequest
	if err := h.readJSON(w, r, &request); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(request.UserID != uuid.Nil, "user_id", "must be provided")
	v.Check(request.OrderID != uuid.Nil, "order_id", "must be provided")
	v.Check(request.Amount > 0, "amount", "must be positive")
	v.Check(request.Hours > 0, "hours", "must be positive")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	entry := &models.LedgerEntry{
		UserID:        request.UserID,
		OperationType: models.OpHold,
		Amount:        request.Amount,
		ExternalKey:   key,
	}

	holdId, err := h.service.Hold(r.Context(), entry, request.OrderID, request.Hours)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLedgerDuplicate):
			h.errorResponse(w, r, http.StatusConflict, "idempotency key already used for another operation")
		case errors.Is(err, storage.ErrInsufficientBalance):
			h.errorResponse(w, r, http.StatusConflict, "not enough balance to hold the order")
		case errors.Is(err, storage.ErrOrderAlreadyHeld):
			h.errorResponse(w, r, http.StatusConflict, "this order already held")
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	response := createHoldResponse{
		HoldID: holdId,
	}

	if err := h.writeJSON(w, http.StatusCreated, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// ConfirmHoldHandler handles POST /users/{user_id}/holds/{order_id}/confirm requests to confirm a held order.
// Validates Idempotency-Key header, path parameters, and calls service.ConfirmWithdraw().
// Returns 200 with hold ID on success.
// Error responses: 400 (bad request, hold not found), 422 (validation), 409 (conflict - duplicate key), 500 (server error).
func (h *APIHandler) ConfirmHoldHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		h.badRequestResponse(w, r, errors.New("Idempotency-Key header is required"))
		return
	}

	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	orderId, err := uuid.Parse(r.PathValue("order_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	v.Check(orderId != uuid.Nil, "order_id", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	entry := &models.LedgerEntry{
		UserID:        userId,
		OperationType: models.OpConfirm,
		Amount:        0,
		ExternalKey:   key,
	}

	holdId, err := h.service.ConfirmWithdraw(r.Context(), entry, orderId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLedgerDuplicate):
			h.errorResponse(w, r, http.StatusConflict, "idempotency key already used for another operation")
		case errors.Is(err, storage.ErrHoldNotFound):
			h.badRequestResponse(w, r, err)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	response := confirmHoldResponse{
		HoldID: holdId,
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// CancelHoldHandler handles POST /users/{user_id}/holds/{order_id}/cancel requests to cancel a held order.
// Validates Idempotency-Key header, path parameters, and calls service.CancelHold().
// Returns 200 with hold ID on success.
// Error responses: 400 (bad request, hold not found), 422 (validation), 409 (conflict - duplicate key), 500 (server error).
func (h *APIHandler) CancelHoldHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		h.badRequestResponse(w, r, errors.New("Idempotency-Key header is required"))
		return
	}

	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	orderId, err := uuid.Parse(r.PathValue("order_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	v.Check(orderId != uuid.Nil, "order_id", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	entry := &models.LedgerEntry{
		UserID:        userId,
		OperationType: models.OpCancel,
		Amount:        0,
		ExternalKey:   key,
	}

	holdId, err := h.service.CancelHold(r.Context(), entry, orderId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLedgerDuplicate):
			h.errorResponse(w, r, http.StatusConflict, "idempotency key already used for another operation")
		case errors.Is(err, storage.ErrHoldNotFound):
			h.badRequestResponse(w, r, err)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	response := cancelHoldResponse{
		HoldID: holdId,
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// GetBalanceHandler handles GET /users/{user_id}/balance requests to retrieve available bonus points.
// Validates path parameters and calls service.GetAvailablePoints().
// Returns 200 with the available balance on success.
// Error responses: 400 (bad request), 422 (validation), 500 (server error).
func (h *APIHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	availableBonuses, err := h.service.GetAvailablePoints(r.Context(), userId)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	response := getBalanceResponse{
		Available: availableBonuses,
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// GetExpirationsHandler handles GET /users/{user_id}/balance/expirations requests to retrieve expiring bonus points.
// Validates path parameters and query parameters, and calls service.GetExpiringAvailablePoints().
// Returns 200 with expiring bonus information on success.
// Error responses: 400 (bad request), 422 (validation), 500 (server error).
func (h *APIHandler) GetExpirationsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	q := r.URL.Query()
	days := q.Get("days")

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	v.Check(days != "", "days", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	intDays, err := strconv.Atoi(days)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	expiringBonuses, err := h.service.GetExpiringAvailablePoints(r.Context(), userId, intDays)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	response := getExpirationsResponse{
		Expiring:  expiringBonuses,
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// GetHistoryHandler handles GET /users/{user_id}/history requests to retrieve transaction history.
// Validates path parameters and pagination query parameters, and calls service.GetHistory().
// Returns 200 with paginated transaction history on success.
// Error responses: 400 (bad request), 422 (validation), 500 (server error).
func (h *APIHandler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	q := r.URL.Query()
	limit := q.Get("limit")
	offset := q.Get("offset")

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	v.Check(limit != "", "limit", "must be provided")
	v.Check(offset != "", "offset", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	history, err := h.service.GetHistory(r.Context(), userId, intLimit, intOffset)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	transactions := make([]transaction, len(history))
	response := GetHistoryResponse{
		Transactions: transactions,
	}

	for i, ledgerEntry := range history {
		transactions[i].OperationType = ledgerEntry.OperationType
		transactions[i].Amount = ledgerEntry.Amount
		transactions[i].CreatedAt = ledgerEntry.CreatedAt
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// GetHoldsHandler handles GET /users/{user_id}/holds requests to retrieve the total amount of held bonus points.
// Validates path parameters and calls service.GetHeld().
// Returns 200 with the held bonus amount on success.
// Error responses: 400 (bad request), 422 (validation), 500 (server error).
func (h *APIHandler) GetHoldsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(userId != uuid.Nil, "user_id", "must be provided")
	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	heldPoints, err := h.service.GetHeld(r.Context(), userId)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	response := GetHoldResponse{
		Amount: heldPoints,
	}

	if err := h.writeJSON(w, http.StatusOK, response, nil); err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
