package service

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"bonus-service/internal/storage"
)

type Worker struct {
	Name     string
	Interval time.Duration
	Job      func(context.Context) error
}

type LeaderService struct {
	repo     *storage.LeaderRepo
	role     string
	instance string
	ttl      int
	isLeader atomic.Bool

	mu      sync.Mutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	workers []Worker

	stopOnce sync.Once
	stopCh   chan struct{}
}

func NewLeaderService(
	repo *storage.LeaderRepo,
	role, instance string,
	ttl int,
) *LeaderService {
	return &LeaderService{
		repo:     repo,
		role:     role,
		instance: instance,
		ttl:      ttl,
		stopCh:   make(chan struct{}),
		workers:  make([]Worker, 0),
	}
}

func (s *LeaderService) RegisterWorker(name string, interval time.Duration, job func(context.Context) error) {
	s.workers = append(s.workers, Worker{
		Name:     name,
		Interval: interval,
		Job:      job,
	})
}

func (s *LeaderService) Start(ctx context.Context) {
	go s.run(ctx)
}

func (s *LeaderService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		s.stopWorkers()
	})
}

func (s *LeaderService) IsLeader() bool {
	return s.isLeader.Load()
}

func (s *LeaderService) run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.stopWorkers()
			return
		case <-s.stopCh:
			s.stopWorkers()
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *LeaderService) tick(ctx context.Context) {
	acquired, err := s.repo.TryAcquireLock(ctx, s.role, s.instance, s.ttl)
	if err != nil {
		log.Printf("[Leader] Error acquiring lock: %v", err)
		return
	}

	if acquired {
		if !s.isLeader.Load() {
			s.isLeader.Store(true)
			log.Printf("[Leader] %s is now the leader", s.instance)
			s.startWorkers(ctx)
		}
	} else {
		if s.isLeader.Load() {
			s.isLeader.Store(false)
			log.Printf("[Leader] %s lost leadership", s.instance)
			s.stopWorkers()
		}
	}
}

func (s *LeaderService) startWorkers(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	log.Printf("[Leader] Starting %d workers...", len(s.workers))

	workerCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.running = true

	for _, w := range s.workers {
		s.wg.Add(1)
		go s.runWorker(workerCtx, w)
	}
}

func (s *LeaderService) stopWorkers() {
	s.mu.Lock()

	if !s.running {
		s.mu.Unlock()
		return
	}

	cancel := s.cancel
	s.cancel = nil
	s.running = false

	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	s.wg.Wait()
}

func (s *LeaderService) runWorker(ctx context.Context, w Worker) {
	defer s.wg.Done()

	log.Printf("[%s] Started", w.Name)

	select {
	case <-ctx.Done():
		return
	default:
	}

	if err := w.Job(ctx); err != nil {
		log.Printf("[%s] Error: %v", w.Name, err)
	}

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Stopped", w.Name)
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := w.Job(ctx); err != nil {
				log.Printf("[%s] Error: %v", w.Name, err)
			}
		}
	}
}
