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
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func getWorkerInterval(envKey string, defaultMinutes int) time.Duration {
	val := os.Getenv(envKey)
	if val == "" {

		return time.Duration(defaultMinutes) * time.Minute
	}
	minutes, err := strconv.Atoi(val)

	if err != nil {
		log.Printf("Invalid %s: %s, using default %d minutes", envKey, val, defaultMinutes)
		return time.Duration(defaultMinutes) * time.Minute
	}
	if minutes <= 0 {
		log.Printf("%s must be positive (%d), using default %d minutes", envKey, minutes, defaultMinutes)
		return time.Duration(defaultMinutes) * time.Minute
	}
	return time.Duration(minutes) * time.Minute
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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

	leaderService.RegisterWorker(
		"TTLWorker",
		getWorkerInterval("WORKER_TTL_INTERVAL", 5),
		bonusService.TTLWorker,
	)
	leaderService.RegisterWorker(
		"BatchExpiryWorker",
		getWorkerInterval("WORKER_BATCH_EXPIRY_INTERVAL", 5),
		bonusService.BatchExpiryWorker,
	)
	leaderService.RegisterWorker(
		"BatchCleanupWorker",
		getWorkerInterval("WORKER_BATCH_CLEANUP_INTERVAL", 60),
		bonusService.BatchCleanupWorker,
	)
	leaderService.RegisterWorker(
		"HoldCleanupWorker",
		getWorkerInterval("WORKER_HOLD_CLEANUP_INTERVAL", 60),
		bonusService.HoldCleanupWorker,
	)

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
