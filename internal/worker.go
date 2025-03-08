package internal

import (
	"log"
	"sync"
)

// WorkerPool is a module that creates workers to fetch metrics and export them in JSON files.
type WorkerPool struct {
	input chan string

	url      string
	from     string
	to       string
	interval string

	wg *sync.WaitGroup
}

// NewWorkerPool returns a WorkerPool instance.
func NewWorkerPool(url, from, to, interval string, poolSize, totalInput int) *WorkerPool {
	instance := WorkerPool{
		url:      url,
		from:     from,
		to:       to,
		interval: interval,
		input:    make(chan string),
		wg:       &sync.WaitGroup{},
	}

	// start workers
	for range poolSize {
		go instance.startNewWorker()
	}

	// set waitgroup
	instance.wg.Add(totalInput)

	return &instance
}

// Collect sends a metric to worker pool.
func (w *WorkerPool) Collect(metric string) {
	w.input <- metric
}

// StopAndWait for all workers to finish.
func (w *WorkerPool) StopAndWait() {
	w.wg.Wait()
}

// startNewWorker creates a process that fetches metrics from Prometheus and stores them in JSON files.
func (w *WorkerPool) startNewWorker() {
	for {
		// get metric from input channel
		metric := <-w.input

		// create HTTP GET request
		req, err := newHttpGetRequest(w.url)
		if err != nil {
			log.Printf("[ERR] build HTTP request failed: %v\n", err)
			w.wg.Done()
			continue
		}

		// set query parameters
		q := req.URL.Query()
		q.Add("query", metric)
		q.Add("start", w.from)
		q.Add("end", w.to)
		q.Add("step", w.interval)
		req.URL.RawQuery = q.Encode()

		// set the headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// fetch metrics
		_, err = fetchMetrics(req)
		if err != nil {
			log.Printf("[ERR] fetch metrics of %s failed: %v\n", metric, err)
			w.wg.Done()
			continue
		}
	}
}
