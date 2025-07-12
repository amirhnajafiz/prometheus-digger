package worker

import (
	"fmt"
	"sync"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/internal/logger"
	"github.com/amirhnajafiz/prometheus-digger/internal/models"
	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// WorkerPool is a module that creates workers to fetch metrics and export them in JSON files.
type WorkerPool struct {
	url     string
	from    string
	to      string
	queries []models.Query
	wg      *sync.WaitGroup
	input   chan *models.Query
}

// NewWorkerPool returns a WorkerPool instance.
func NewWorkerPool(cfg *configs.Config) *WorkerPool {
	// create worker pool
	instance := WorkerPool{
		url:     fmt.Sprintf("%s%s", cfg.URL, promAPI),
		from:    cfg.From,
		to:      cfg.To,
		queries: cfg.Queries,
		input:   make(chan *models.Query),
		wg:      &sync.WaitGroup{},
	}

	// check if the output directory exists
	if err := pkg.CheckDir(outputDir); err != nil {
		logger.Error(fmt.Sprintf("check output directory failed: %v", err))
		return nil
	}

	// start workers
	poolSize := min(len(cfg.Queries)/3, 1)
	for range poolSize {
		go instance.startNewWorker()
	}

	return &instance
}

// Collect sends a metric to worker pool.
func (w *WorkerPool) Collect() {
	for _, q := range w.queries {
		w.wg.Add(1)
		w.input <- &q
	}
}

// StopAndWait for all workers to finish.
func (w *WorkerPool) StopAndWait() {
	w.wg.Wait()
}
