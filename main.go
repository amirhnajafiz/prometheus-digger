package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/amirhnajafiz/prometheus-digger/internal"
)

var (
	prometheusUrl string
	from          string
	to            string
	interval      string
	metrics       []string
)

// load the vars from environment variables if not set
func initVars() {
	if prometheusUrl == "" {
		prometheusUrl = os.Getenv("PD_PROMETHEUS_URL")
		if prometheusUrl == "" {
			prometheusUrl = "http://localhost:9090"
		}
	}

	if len(metrics) == 0 {
		metricsEnv := os.Getenv("PD_METRICS")
		if metricsEnv == "" {
			metrics = []string{"http_requests_total"}
		} else {
			metrics = strings.Split(metricsEnv, ",")
		}
	}

	if from == "" {
		from = os.Getenv("PD_FROM")
		if from == "" {
			from = "2023-01-01T00:00:00Z"
		}
	}

	if to == "" {
		to = os.Getenv("PD_TO")
		if to == "" {
			to = "2023-01-01T00:00:00Z"
		}
	}

	if interval == "" {
		interval = os.Getenv("PD_INTERVAL")
		if interval == "" {
			interval = "1m"
		}
	}
}

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

	// initialize variables
	initVars()

	// loop through the metrics and create a goroutine for each
	for _, metric := range metrics {
		go func(metric string) {
			// create a new worker
			resp, err := internal.Worker(prometheusUrl, metric, from, to, interval)
			if err != nil {
				panic(err)
			}

			// print the response
			fmt.Println(resp)
		}(metric)
	}
}
