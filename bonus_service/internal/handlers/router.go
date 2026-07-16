package handlers

import "net/http"

// NewRouter creates and configures the HTTP router with all API endpoints.
func NewRouter(h *APIHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /accrual", h.AccrualHandler)
	router.HandleFunc("POST /hold", h.CreateHoldHandler)
	router.HandleFunc("POST /users/{user_id}/holds/{order_id}/confirm", h.ConfirmHoldHandler)
	router.HandleFunc("POST /users/{user_id}/holds/{order_id}/cancel", h.CancelHoldHandler)
	router.HandleFunc("GET /users/{user_id}/balance", h.GetBalanceHandler)
	router.HandleFunc("GET /users/{user_id}/balance/expirations", h.GetExpirationsHandler)
	router.HandleFunc("GET /users/{user_id}/history", h.GetHistoryHandler)
	router.HandleFunc("GET /users/{user_id}/holds", h.GetHoldsHandler)

	return router
}