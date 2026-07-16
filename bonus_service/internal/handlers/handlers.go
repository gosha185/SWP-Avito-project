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

type APIHandler struct {
	service *service.BonusService
}

func NewAPIHandler(s *service.BonusService) *APIHandler {
	return &APIHandler{service: s}
}

//----------------------------- HANDLERS ------------------------------------------

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
