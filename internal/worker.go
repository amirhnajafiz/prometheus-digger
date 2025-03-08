package internal

import "log"

type WorkerPool struct {
	input chan string

	url      string
	from     string
	to       string
	interval string
}

func NewWorkerPool(url, from, to, interval string, poolSize int) *WorkerPool {
	return &WorkerPool{
		url:      url,
		from:     from,
		to:       to,
		interval: interval,
		input:    make(chan string),
	}
}

func (w *WorkerPool) Collect(metric string) {
	w.input <- metric
}

// worker fetches metrics from Prometheus.
// It takes the Prometheus URL, metric name, start time, end time, and interval as parameters.
func (w *WorkerPool) startNewWorker() {
	for {
		// get metric from input channel
		metric := <-w.input

		// create HTTP GET request
		req, err := newHttpGetRequest(w.url)
		if err != nil {
			log.Printf("[ERR] build HTTP request failed: %v\n", err)
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
			continue
		}
	}
}
