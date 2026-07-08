package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"testing"
	"time"

	"bonus-service/internal/models"
	"bonus-service/internal/service"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const USER_COUNT = 999
const OPERATIONS = 200

const fullSchema = `
CREATE TABLE IF NOT EXISTS balances (
    user_id    UUID PRIMARY KEY,
    available  BIGINT      NOT NULL DEFAULT 0 CHECK (available >= 0),
    held       BIGINT      NOT NULL DEFAULT 0 CHECK (held >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS batches (
    id         BIGSERIAL,
    user_id    UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    remaining  BIGINT      NOT NULL CHECK (remaining >= 0),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, id)
);

CREATE INDEX IF NOT EXISTS idx_batches_user_expires ON batches (user_id, expires_at) WHERE remaining > 0;
CREATE INDEX IF NOT EXISTS idx_batches_expired ON batches (expires_at) WHERE remaining > 0;

CREATE TABLE IF NOT EXISTS holds (
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    order_id   UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    status     VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'confirmed', 'cancelled')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, order_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_holds_user_id_id ON holds (user_id, id);

CREATE TABLE IF NOT EXISTS hold_batches (
    hold_id       BIGINT NOT NULL REFERENCES holds (id) ON DELETE CASCADE,
    batch_user_id UUID   NOT NULL,
    batch_id      BIGINT NOT NULL,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    PRIMARY KEY (hold_id, batch_user_id, batch_id),
    FOREIGN KEY (batch_user_id, batch_id) REFERENCES batches (user_id, id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_holds_expires_status ON holds (expires_at, status) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_holds_order_id ON holds (order_id);
CREATE INDEX IF NOT EXISTS idx_holds_user_status ON holds (user_id, status);
CREATE INDEX IF NOT EXISTS idx_hold_batches_batch_id ON hold_batches (batch_user_id, batch_id);

CREATE TABLE IF NOT EXISTS ledger (
    id             BIGSERIAL,
    user_id        UUID         NOT NULL,
    operation_type VARCHAR(20)  NOT NULL CHECK (operation_type IN ('accrual','hold','confirm','cancel','expiry')),
    amount         BIGINT       NOT NULL CHECK (amount > 0),
    batch_id       BIGINT,
    external_key   VARCHAR(255) NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    metadata       JSONB,
    PRIMARY KEY (user_id, id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ledger_external_key ON ledger (external_key);
CREATE INDEX IF NOT EXISTS idx_ledger_user_created ON ledger (user_id, created_at);
`

type IntegrationPerfResult struct {
	Phase        string
	Operation    string
	Count        int
	TotalTime    time.Duration
	MinTime      time.Duration
	MaxTime      time.Duration
	AvgTime      time.Duration
	P50Time      time.Duration
	P95Time      time.Duration
	P99Time      time.Duration
	OpsPerSecond float64
}

func printIntegrationResults(results []IntegrationPerfResult) {
	fmt.Println("\n=== Integration Performance Test ===")

	for _, r := range results {
		if r.Phase != "" {
			fmt.Printf("\n[%s]\n", r.Phase)
		}
		fmt.Printf("%-45s : count=%d\n", r.Operation, r.Count)
		fmt.Printf("  Total: %v\n", r.TotalTime.Round(time.Millisecond))
		fmt.Printf("  Min:   %v\n", r.MinTime.Round(time.Microsecond))
		fmt.Printf("  Max:   %v\n", r.MaxTime.Round(time.Microsecond))
		fmt.Printf("  Avg:   %v\n", r.AvgTime.Round(time.Microsecond))
		fmt.Printf("  P50:   %v\n", r.P50Time.Round(time.Microsecond))
		fmt.Printf("  P95:   %v\n", r.P95Time.Round(time.Microsecond))
		fmt.Printf("  P99:   %v\n", r.P99Time.Round(time.Microsecond))
		fmt.Printf("  RPS:   %.2f ops/sec\n", r.OpsPerSecond)
	}
}

func calculatePercentile(times []time.Duration, percentile float64) time.Duration {
	if len(times) == 0 {
		return 0
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	index := int(float64(len(times)-1) * percentile)
	return times[index]
}

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

	if _, err = db.Exec(fullSchema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	return db
}

func setupTestService(t *testing.T) *service.BonusService {
	db := newTestDB(t)
	return service.NewBonusService(db)
}

func TestService_IntegrationPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration performance test in short mode")
	}

	svc := setupTestService(t)
	ctx := context.Background()
	results := []IntegrationPerfResult{}

	const userCount = USER_COUNT
	const operationCount = OPERATIONS

	fmt.Println("Starting integration performance test...")
	fmt.Printf("Users: %d, Operations per type: %d\n\n", userCount, operationCount)

	// SETUP
	fmt.Println("Phase 1: Creating users with balance...")
	userIDs := make([]uuid.UUID, userCount)
	times := make([]time.Duration, 0, userCount)

	start := time.Now()
	for i := 0; i < userCount; i++ {
		userID := uuid.New()
		userIDs[i] = userID

		entry := &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpAccrual,
			Amount:        5000,
			ExternalKey:   uuid.New().String(),
		}

		opStart := time.Now()
		_, err := svc.Accrue(ctx, entry, int64(i%30+1))
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to accrue for user %d: %v", i, err)
		}
	}
	totalTime := time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "SETUP",
		Operation:    "Create users with balance",
		Count:        userCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(userCount) / totalTime.Seconds(),
	})

	// ACCRUE
	fmt.Println("\nPhase 2: Testing Accrue operations...")
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		entry := &models.LedgerEntry{
			UserID:        userIDs[i%userCount],
			OperationType: models.OpAccrual,
			Amount:        100,
			ExternalKey:   uuid.New().String(),
		}

		opStart := time.Now()
		_, err := svc.Accrue(ctx, entry, 30)
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to accrue: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "ACCRUE",
		Operation:    "Accrue additional points",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// HOLD
	fmt.Println("\nPhase 3: Testing Hold operations...")
	orderIDs := make([]uuid.UUID, operationCount)
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		entry := &models.LedgerEntry{
			UserID:        userIDs[i%userCount],
			OperationType: models.OpHold,
			Amount:        200,
			ExternalKey:   uuid.New().String(),
		}
		orderID := uuid.New()
		orderIDs[i] = orderID

		opStart := time.Now()
		_, err := svc.Hold(ctx, entry, orderID, 24)
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to hold: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "HOLD",
		Operation:    "Hold points for orders",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// QUERY: GetAvailablePoints
	fmt.Println("\nPhase 4: Testing GetAvailablePoints...")
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		opStart := time.Now()
		_, err := svc.GetAvailablePoints(ctx, userIDs[i%userCount])
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to get available points: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "QUERY",
		Operation:    "GetAvailablePoints",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// QUERY: GetExpiringAvailablePoints
	fmt.Println("\nPhase 5: Testing GetExpiringAvailablePoints...")
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		opStart := time.Now()
		_, err := svc.GetExpiringAvailablePoints(ctx, userIDs[i%userCount], 30)
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to get expiring points: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "QUERY",
		Operation:    "GetExpiringAvailablePoints",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// QUERY: GetHeld
	fmt.Println("\nPhase 6: Testing GetHeld...")
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		opStart := time.Now()
		_, err := svc.GetHeld(ctx, userIDs[i%userCount])
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to get held: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "QUERY",
		Operation:    "GetHeld",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// QUERY: GetHistory
	fmt.Println("\nPhase 7: Testing GetHistory...")
	times = make([]time.Duration, 0, operationCount)

	start = time.Now()
	for i := 0; i < operationCount; i++ {
		opStart := time.Now()
		_, err := svc.GetHistory(ctx, userIDs[i%userCount], 50, 0)
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to get history: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "QUERY",
		Operation:    "GetHistory",
		Count:        operationCount,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount) / totalTime.Seconds(),
	})

	// CONFIRM
	fmt.Println("\nPhase 8: Testing ConfirmWithdraw...")
	times = make([]time.Duration, 0, operationCount/2)

	start = time.Now()
	for i := 0; i < operationCount/2; i++ {
		entry := &models.LedgerEntry{
			UserID:        userIDs[i%userCount],
			OperationType: models.OpConfirm,
			ExternalKey:   uuid.New().String(),
		}

		opStart := time.Now()
		_, err := svc.ConfirmWithdraw(ctx, entry, orderIDs[i])
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to confirm: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "CONFIRM",
		Operation:    "ConfirmWithdraw",
		Count:        operationCount / 2,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount/2) / totalTime.Seconds(),
	})

	// CANCEL
	fmt.Println("\nPhase 9: Testing CancelHold...")
	times = make([]time.Duration, 0, operationCount/2)

	start = time.Now()
	for i := operationCount / 2; i < operationCount; i++ {
		entry := &models.LedgerEntry{
			UserID:        userIDs[i%userCount],
			OperationType: models.OpCancel,
			ExternalKey:   uuid.New().String(),
		}

		opStart := time.Now()
		_, err := svc.CancelHold(ctx, entry, orderIDs[i])
		times = append(times, time.Since(opStart))

		if err != nil {
			t.Fatalf("failed to cancel: %v", err)
		}
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "CANCEL",
		Operation:    "CancelHold",
		Count:        operationCount / 2,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount/2) / totalTime.Seconds(),
	})

	// FULL SCENARIO
	fmt.Println("\nPhase 10: Testing full scenario (Accrue->Hold->Confirm)...")
	times = make([]time.Duration, 0, operationCount/2)

	start = time.Now()
	for i := 0; i < operationCount/2; i++ {
		userID := uuid.New()
		orderID := uuid.New()

		scenarioStart := time.Now()

		entry := &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpAccrual,
			Amount:        1000,
			ExternalKey:   uuid.New().String(),
		}
		_, err := svc.Accrue(ctx, entry, 30)
		if err != nil {
			t.Fatalf("failed to accrue: %v", err)
		}

		entry = &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpHold,
			Amount:        500,
			ExternalKey:   uuid.New().String(),
		}
		_, err = svc.Hold(ctx, entry, orderID, 24)
		if err != nil {
			t.Fatalf("failed to hold: %v", err)
		}

		entry = &models.LedgerEntry{
			UserID:        userID,
			OperationType: models.OpConfirm,
			ExternalKey:   uuid.New().String(),
		}
		_, err = svc.ConfirmWithdraw(ctx, entry, orderID)
		if err != nil {
			t.Fatalf("failed to confirm: %v", err)
		}

		times = append(times, time.Since(scenarioStart))
	}
	totalTime = time.Since(start)

	results = append(results, IntegrationPerfResult{
		Phase:        "FULL SCENARIO",
		Operation:    "Accrue -> Hold -> Confirm",
		Count:        operationCount / 2,
		TotalTime:    totalTime,
		MinTime:      calculatePercentile(times, 0.0),
		MaxTime:      calculatePercentile(times, 1.0),
		AvgTime:      calculatePercentile(times, 0.5),
		P50Time:      calculatePercentile(times, 0.5),
		P95Time:      calculatePercentile(times, 0.95),
		P99Time:      calculatePercentile(times, 0.99),
		OpsPerSecond: float64(operationCount/2) / totalTime.Seconds(),
	})

	printIntegrationResults(results)
}
