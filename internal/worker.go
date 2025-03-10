package internal

import (
	"fmt"
	"log"
	"sync"

	"github.com/amirhnajafiz/prometheus-digger/pkg"
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

	// check if the output directory exists
	err := pkg.CheckDir(outputDir)
	if err != nil {
		log.Printf("[ERR] check output directory failed: %v\n", err)
		return nil
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
		req, err := pkg.NewHttpGetRequest(w.url)
		if err != nil {
			w.throwError(fmt.Sprintf("build HTTP request for %s failed: %v", metric, err))
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
		resp, err := pkg.FetchMetrics(req)
		if err != nil {
			w.throwError(fmt.Sprintf("fetch metrics of %s failed: %v", metric, err))
			continue
		}

		// check the output directory
		if err := pkg.CheckDir(outputDir + "/" + metric); err != nil {
			w.throwError(fmt.Sprintf("check output directory of %s failed: %v", metric, err))
			continue
		}

		// store metrics in JSON file
		if err = pkg.WriteFile(w.getFileName(metric), resp); err != nil {
			w.throwError(fmt.Sprintf("store metrics of %s failed: %v", metric, err))
			continue
		}

		log.Printf("[INFO] metrics of %s stored successfully.\n", metric)
		w.wg.Done()
	}
}

// getFileName returns the file name for the given metric, from and to.
func (w *WorkerPool) getFileName(metric string) string {
	return outputDir + "/" + metric + "/" + w.from + "_" + w.to + ".json"
}

// throwError logs an error message and marks the worker as done.
func (w *WorkerPool) throwError(msg string) {
	log.Printf("[ERR] %s\n", msg)
	w.wg.Done()
}
