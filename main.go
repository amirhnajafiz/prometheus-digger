package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal"
	"github.com/amirhnajafiz/prometheus-digger/pkg"
)

const (
	// this const value is the number of workers per metrics, ex. 1 worker per 3 metric
	poolSizeMetricsRatio = 3
	// this const value is the path of the prometheus API
	promAPI = "/api/v1/query_range"
)

var (
	// program flags
	prometheusUrl string
	from          string
	to            string
	interval      string
	metricsFile   string
	metrics       []string
)

func main() {
	// define flags for the vars
	flag.StringVar(&prometheusUrl, "prometheus-url", "http://127.0.0.1:9090", "Prometheus URL")
	flag.StringVar(&from, "from", time.Now().String(), "Start time for the query")
	flag.StringVar(&to, "to", time.Now().Add(1*time.Hour).String(), "End time for the query")
	flag.StringVar(&interval, "interval", "1m", "Interval for the query")
	flag.StringVar(&metricsFile, "metrics-file", "", "Path to the metrics file")
	metricsFlag := flag.String("metrics", "", "Metrics to query")

	// parse the flags
	flag.Parse()

	// split metrics flag if set
	if *metricsFlag != "" {
		metrics = strings.Split(*metricsFlag, ",")
	}

	// read metrics from file if set
	if metricsFile != "" {
		bytes, err := pkg.ReadFile(metricsFile)
		if err != nil {
			fmt.Println("[ERR] failed to read metrics file:", err)
			return
		}

		jsons, err := internal.BytesToJSONs(bytes)
		if err != nil {
			fmt.Println("[ERR] failed to parse metrics file:", err)
			return
		}
	}

	// create a worker pool
	total := len(metrics)
	poolSize := total/poolSizeMetricsRatio + 1
	pool := internal.NewWorkerPool(
		fmt.Sprintf("%s%s", prometheusUrl, promAPI),
		from,
		to,
		interval,
		poolSize,
		total,
	)
	if pool == nil {
		fmt.Println("[ERR] failed to create worker pool!")
		return
	}

	// loop through the metrics and create a goroutine for each
	for _, metric := range metrics {
		pool.Collect(metric)
	}

	// wait for all workers
	pool.StopAndWait()
}
