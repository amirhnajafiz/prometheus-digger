package worker

import (
	"fmt"

	"github.com/amirhnajafiz/prometheus-digger/internal/config"
	"github.com/amirhnajafiz/prometheus-digger/internal/logger"
	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

// startNewWorker creates a process that fetches metrics from Prometheus and stores them in JSON files.
func (w *WorkerPool) startNewWorker() {
	for {
		// get metric from input channel
		query := <-w.input

		// set the callback function
		var callback func(*config.Query) ([]byte, error)
		if len(query.Metric) < 32 {
			callback = w.followGET
		} else {
			callback = w.followPOST
		}

		// fetch metrics
		resp, err := callback(query)
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

		logger.Info(fmt.Sprintf("metrics of %s stored successfully", query.Name))
		w.wg.Done()
	}
}

// followGET sends an HTTP GET request to the Prometheus server and returns the response body.
func (w *WorkerPool) followGET(query *config.Query) ([]byte, error) {
	logger.Info(fmt.Sprintf("metrics of %s are being pulled by GET", query.Name))

	// create HTTP GET request
	req, err := pkg.NewHttpGetRequest(w.cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("build HTTP request failed: %v", err)
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
	return pkg.FetchMetrics(req)
}

// followPOST sends an HTTP POST request to the Prometheus server and returns the response body.
func (w *WorkerPool) followPOST(query *config.Query) ([]byte, error) {
	logger.Info(fmt.Sprintf("metrics of %s are being pulled by POST", query.Name))

	// set the body
	body := fmt.Sprintf(
		"query=%s&start=%s&end=%s&step=%s",
		query.Metric,
		w.cfg.From,
		w.cfg.To,
		query.Interval,
	)

	// create HTTP POST request
	req, err := pkg.NewHttpPostRequest(w.cfg.URL, []byte(body))
	if err != nil {
		return nil, fmt.Errorf("build HTTP request failed: %v", err)
	}

	// set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// fetch metrics
	return pkg.FetchMetrics(req)
}

// getFileName returns the file name for the given metric, from and to.
func (w *WorkerPool) getFileName(metric string) string {
	return outputDir + "/" + metric + "/" + w.cfg.From + "_" + w.cfg.To + ".json"
}

// throwError logs an error message and marks the worker as done.
func (w *WorkerPool) throwError(msg string) {
	logger.Error(msg)
	w.wg.Done()
}
