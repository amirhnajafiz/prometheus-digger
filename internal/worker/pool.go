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
	url       string
	from      string
	to        string
	queries   []models.Query
	wg        *sync.WaitGroup
	input     chan *models.Query
	latencies []float64
}

// NewWorkerPool returns a WorkerPool instance.
func NewWorkerPool(cfg *configs.Config) *WorkerPool {
	// create worker pool
	instance := WorkerPool{
		url:       fmt.Sprintf("%s%s", cfg.URL, promAPI),
		from:      cfg.From,
		to:        cfg.To,
		queries:   cfg.Queries,
		input:     make(chan *models.Query),
		wg:        &sync.WaitGroup{},
		latencies: make([]float64, 0),
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

// PrintStats prints the latencies of the workers.
func (w *WorkerPool) PrintStats() {
	if len(w.latencies) == 0 {
		return
	}

	fmt.Println("")

	// calculate average latency
	var totalLatency float64
	for _, latency := range w.latencies {
		totalLatency += latency
	}
	avgLatency := totalLatency / float64(len(w.latencies))

	fmt.Println("|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|")
	fmt.Printf("| %-30s | %-30s | %-30s | %-30s | %-30s |\n", "Target", "From", "To", "Requests", "Latency (ms)")
	fmt.Println("|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|")

	// split long text into multiple lines if needed
	target := w.url
	targetLines := []string{}
	for len(target) > 30 {
		targetLines = append(targetLines, target[:30])
		target = target[30:]
	}
	if len(target) > 0 {
		targetLines = append(targetLines, target)
	}

	from := w.from
	fromLines := []string{}
	for len(from) > 30 {
		fromLines = append(fromLines, from[:30])
		from = from[30:]
	}
	if len(from) > 0 {
		fromLines = append(fromLines, from)
	}

	to := w.to
	toLines := []string{}
	for len(to) > 30 {
		toLines = append(toLines, to[:30])
		to = to[30:]
	}
	if len(to) > 0 {
		toLines = append(toLines, to)
	}

	// find the maximum number of lines needed
	maxLines := len(targetLines)
	if len(fromLines) > maxLines {
		maxLines = len(fromLines)
	}
	if len(toLines) > maxLines {
		maxLines = len(toLines)
	}

	// print each line
	for i := 0; i < maxLines; i++ {
		targetPart := ""
		if i < len(targetLines) {
			targetPart = targetLines[i]
		}

		fromPart := ""
		if i < len(fromLines) {
			fromPart = fromLines[i]
		}

		toPart := ""
		if i < len(toLines) {
			toPart = toLines[i]
		}

		if i == 0 {
			// First line includes requests and latency
			fmt.Printf("| %-30s | %-30s | %-30s | %-30d | %-30.2f |\n", targetPart, fromPart, toPart, len(w.latencies), avgLatency)
		} else {
			// Subsequent lines only show continuation text
			fmt.Printf("| %-30s | %-30s | %-30s | %-30s | %-30s |\n", targetPart, fromPart, toPart, "", "")
		}
	}

	fmt.Println("|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|")
}
