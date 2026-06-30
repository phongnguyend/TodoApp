package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/database"
	"github.com/todo/backend/go/internal/worker/job"
)

func main() {
	// ── Configuration ─────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("[worker] failed to connect to database: %v", err)
	}

	interval := time.Duration(cfg.WorkerIntervalMinutes) * time.Minute
	log.Printf("[worker] starting background worker (interval=%v)", interval)

	// Run once immediately on startup — mirrors the Python worker behaviour.
	job.SendIncompleteTodosEmail(db, cfg)

	// ── Ticker loop ───────────────────────────────────────────────────────────
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ticker.C:
			job.SendIncompleteTodosEmail(db, cfg)
		case sig := <-quit:
			log.Printf("[worker] received signal %s — shutting down", sig)
			return
		}
	}
}
