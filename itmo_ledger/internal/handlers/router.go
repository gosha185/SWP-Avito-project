package handlers

import "net/http"

func NewRouter(h *APIHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /v1/accrual", h.AccrualHandler)
	router.HandleFunc("POST /v1/hold", h.CreateHoldHandler)
	router.HandleFunc("POST /v1/users/{user_id}/holds/{order_id}/confirm", h.ConfirmHoldHandler)
	router.HandleFunc("POST /v1/users/{user_id}/holds/{order_id}/cancel", h.CancelHoldHandler)
	router.HandleFunc("GET /v1/users/{user_id}/balance", h.GetBalanceHandler)
	router.HandleFunc("GET /v1/users/{user_id}/history", h.GetHistoryHandler)

	return router
}