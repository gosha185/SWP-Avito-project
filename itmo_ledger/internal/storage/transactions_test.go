package storage

import (
	"context"
	"database/sql"
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

// applyTransaction replicates the read-modify-write logic from createTransactionHandler.
func applyTransaction(model data.BalanceModel, userID uuid.UUID, amount int, txType string) error {
	balance, err := model.Get(userID)
	if err != nil {
		if err == data.ErrRecordNotFound {
			return model.Insert(&data.Balance{Id: userID, Amount: amount})
		}
		return err
	}

	if txType == "deposit" {
		balance.Amount += amount
	} else {
		balance.Amount -= amount
	}
	return model.Update(balance)
}

// TestConcurrentDepositsLoseUpdates demonstrates the lost-update race condition.
//
// 10 goroutines each deposit 10 points simultaneously.
// Expected final balance: initialDeposit + goroutines*depositAmount.
func TestConcurrentDepositsLoseUpdates(t *testing.T) {
	db := newTestDB(t)
	model := data.BalanceModel{DB: db}

	userID := uuid.New()
	initialDeposit := 1
	if err := applyTransaction(model, userID, initialDeposit, "deposit"); err != nil {
		t.Fatalf("seed deposit: %v", err)
	}

	goroutines := 10
	depositAmount := 10

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = applyTransaction(model, userID, depositAmount, "deposit")
		}()
	}
	wg.Wait()

	balance, err := model.Get(userID)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}

	want := initialDeposit + goroutines*depositAmount
	got := balance.Amount

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
	model := data.BalanceModel{DB: db}

	userID := uuid.New()
	seed := 1000
	if err := applyTransaction(model, userID, seed, "deposit"); err != nil {
		t.Fatalf("seed deposit: %v", err)
	}

	deposits := 20
	withdrawals := 10
	depositAmt := 10
	withdrawAmt := 5

	// Serial expectation: seed + deposits*depositAmt - withdrawals*withdrawAmt
	want := seed + deposits*depositAmt - withdrawals*withdrawAmt

	var wg sync.WaitGroup
	for i := 0; i < deposits; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = applyTransaction(model, userID, depositAmt, "deposit")
		}()
	}
	for i := 0; i < withdrawals; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = applyTransaction(model, userID, withdrawAmt, "withdrawal")
		}()
	}
	wg.Wait()

	balance, err := model.Get(userID)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}

	got := balance.Amount
	if got != want {
		t.Errorf("RACE CONDITION: mixed concurrent transactions diverged — want %d, got %d (delta %+d)",
			want, got, got-want)
	}
}