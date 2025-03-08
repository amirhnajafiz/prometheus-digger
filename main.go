package main

import (
	"flag"
	"strings"

	"github.com/amirhnajafiz/prometheus-digger/internal"
)

const (
	// this const value is the number of workers per metrics, ex. 1 worker per 3 metric
	poolSizeMetricsRatio = 3
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
	pool := internal.NewWorkerPool(
		prometheusUrl,
		from,
		to,
		interval,
		len(metrics)/poolSizeMetricsRatio,
	)

	// loop through the metrics and create a goroutine for each
	for _, metric := range metrics {
		pool.Collect(metric)
	}
}
