package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/amirhnajafiz/prometheus-digger/internal"
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
	metrics       []string
)

func main() {
	// define flags for the vars
	flag.StringVar(&prometheusUrl, "prometheus-url", "", "Prometheus URL")
	flag.StringVar(&from, "from", "", "Start time for the query")
	flag.StringVar(&to, "to", "", "End time for the query")
	flag.StringVar(&interval, "interval", "", "Interval for the query")
	metricsFlag := flag.String("metrics", "", "Metrics to query")

	// parse the flags
	flag.Parse()

	// split metrics flag if set
	if *metricsFlag != "" {
		metrics = strings.Split(*metricsFlag, ",")
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
