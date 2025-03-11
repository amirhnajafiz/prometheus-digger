package worker

import (
	"fmt"
	"log"
	"sync"

	"github.com/amirhnajafiz/prometheus-digger/internal/config"
	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// WorkerPool is a module that creates workers to fetch metrics and export them in JSON files.
type WorkerPool struct {
	cfg   *config.Config
	wg    *sync.WaitGroup
	input chan *config.Query
}

// NewWorkerPool returns a WorkerPool instance.
func NewWorkerPool(cfg *config.Config) *WorkerPool {
	// set the Prometheus API
	cfg.URL = fmt.Sprintf("%s%s", cfg.URL, promAPI)

	// create worker pool
	instance := WorkerPool{
		cfg:   cfg,
		input: make(chan *config.Query),
		wg:    &sync.WaitGroup{},
	}

	// check if the output directory exists
	if err := pkg.CheckDir(outputDir); err != nil {
		log.Printf("[ERR] check output directory failed: %v\n", err)
		return nil
	}

	// start workers
	for range cfg.PoolSize {
		go instance.startNewWorker()
	}

	return &instance
}

// Collect sends a metric to worker pool.
func (w *WorkerPool) Collect() {
	for _, q := range w.cfg.Queries {
		w.wg.Add(1)
		w.input <- &q
	}
}

// StopAndWait for all workers to finish.
func (w *WorkerPool) StopAndWait() {
	w.wg.Wait()
}
