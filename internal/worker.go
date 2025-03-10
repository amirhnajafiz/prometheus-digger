package internal

import (
	"fmt"
	"log"

	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// startNewWorker creates a process that fetches metrics from Prometheus and stores them in JSON files.
func (w *WorkerPool) startNewWorker() {
	for {
		// get metric from input channel
		query := <-w.input

		// create HTTP GET request
		req, err := pkg.NewHttpGetRequest(w.cfg.URL)
		if err != nil {
			w.throwError(fmt.Sprintf("build HTTP request for %s failed: %v", query.Name, err))
			continue
		}

		// set query parameters
		q := req.URL.Query()
		q.Add("query", query.Metric)
		q.Add("start", w.cfg.From)
		q.Add("end", w.cfg.To)
		q.Add("step", query.Interval)
		req.URL.RawQuery = q.Encode()

		// set the headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// fetch metrics
		resp, err := pkg.FetchMetrics(req)
		if err != nil {
			w.throwError(fmt.Sprintf("fetch metrics of %s failed: %v", query.Name, err))
			continue
		}

		// check the output directory
		if err := pkg.CheckDir(outputDir + "/" + query.Name); err != nil {
			w.throwError(fmt.Sprintf("check output directory of %s failed: %v", query.Name, err))
			continue
		}

		// store metrics in JSON file
		if err = pkg.WriteFile(w.getFileName(query.Name), resp); err != nil {
			w.throwError(fmt.Sprintf("store metrics of %s failed: %v", query.Name, err))
			continue
		}

		log.Printf("[INFO] metrics of %s stored successfully.\n", query.Name)
		w.wg.Done()
	}
}

// getFileName returns the file name for the given metric, from and to.
func (w *WorkerPool) getFileName(metric string) string {
	return outputDir + "/" + metric + "/" + w.cfg.From + "_" + w.cfg.To + ".json"
}

// throwError logs an error message and marks the worker as done.
func (w *WorkerPool) throwError(msg string) {
	log.Printf("[ERR] %s\n", msg)
	w.wg.Done()
}
