package main

import (
	"bonus-service/internal/handlers"
	"bonus-service/internal/service"
	"bonus-service/internal/storage"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	db, err := storage.NewDB(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	bonusService := service.NewBonusService(db)
	apiHandler := handlers.NewAPIHandler(bonusService)
	router := handlers.NewRouter(apiHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname: ", err)
	}

	leaderService := service.NewLeaderService(
		storage.NewLeaderRepo(db),
		"bonus_worker",
		hostname+":9091",
		30,
	)

	leaderService.RegisterWorker("TTLWorker", 5*time.Minute, bonusService.TTLWorker)
	leaderService.RegisterWorker("BatchExpiryWorker", 5*time.Minute, bonusService.BatchExpiryWorker)
	leaderService.RegisterWorker("BatchCleanupWorker", 1*time.Hour, bonusService.BatchCleanupWorker)
	leaderService.RegisterWorker("HoldCleanupWorker", 1*time.Hour, bonusService.HoldCleanupWorker)

	leaderService.Start(ctx)
	srv := &http.Server{
		Addr:         ":9091",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("HTTP server error:", err)
	}
}
