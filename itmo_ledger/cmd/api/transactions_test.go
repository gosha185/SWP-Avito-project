package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"simple-ledger.itmo.ru/internal/data"
)

// --- test infrastructure ---

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err = db.PingContext(ctx); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS balances (
			id uuid PRIMARY KEY,
			updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
			amount int
		)`)
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return db
}

func newTestApp(db *sql.DB) *application {
	return &application{
		logger: log.New(os.Stderr, "TEST ", 0),
		models: data.NewModels(db),
	}
}

func postTransaction(t *testing.T, app *application, userID uuid.UUID, amount int, txType string) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"user_id": userID.String(),
		"amount":  amount,
		"type":    txType,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	app.createTransactionHandler(rec, req)
	return rec
}

func getBalance(t *testing.T, app *application, userID uuid.UUID) int {
	t.Helper()
	balance, err := app.models.Balances.Get(userID)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}
	return balance.Amount
}

// --- concurrency tests ---

// TestConcurrentDepositsLoseUpdates demonstrates the lost-update race condition.
//
// 10 goroutines each deposit 10 points simultaneously.
// Expected final balance: 100.
func TestConcurrentDepositsLoseUpdates(t *testing.T) {
	db := newTestDB(t)
	app := newTestApp(db)

	userID := uuid.New()
	// Seed a starting balance so all goroutines find an existing record.
	postTransaction(t, app, userID, 0+1, "deposit") // balance = 1 (amount must be > 0)
	// Adjust: we want to start from a known value; subtract 1 via withdrawal.
	// Actually, let's just start with a real deposit and factor it in.
	initialDeposit := 1
	goroutines := 10
	depositAmount := 10

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postTransaction(t, app, userID, depositAmount, "deposit")
		}()
	}
	wg.Wait()

	want := initialDeposit + goroutines*depositAmount
	got := getBalance(t, app, userID)

	if got != want {
		t.Errorf("RACE CONDITION: concurrent deposits lost updates — want %d, got %d (lost %d points)",
			want, got, want-got)
	} else {
		t.Logf("balance correct: %d", got)
	}
}

// TestConcurrentMixedTransactionsBalance runs deposits and withdrawals in
// parallel and verifies the final balance equals what a serial execution
// would produce. Any deviation exposes lost updates.
func TestConcurrentMixedTransactionsBalance(t *testing.T) {
	db := newTestDB(t)
	app := newTestApp(db)

	userID := uuid.New()
	postTransaction(t, app, userID, 1000, "deposit") // seed

	deposits := 20
	withdrawals := 10
	depositAmt := 10
	withdrawAmt := 5

	// Serial expectation: 1000 + 20*10 - 10*5 = 1000+200-50 = 1150
	want := 1000 + deposits*depositAmt - withdrawals*withdrawAmt

	var wg sync.WaitGroup
	for i := 0; i < deposits; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postTransaction(t, app, userID, depositAmt, "deposit")
		}()
	}
	for i := 0; i < withdrawals; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postTransaction(t, app, userID, withdrawAmt, "withdrawal")
		}()
	}
	wg.Wait()

	got := getBalance(t, app, userID)
	if got != want {
		t.Errorf("RACE CONDITION: mixed concurrent transactions diverged — want %d, got %d (delta %+d)",
			want, got, got-want)
	}
}
