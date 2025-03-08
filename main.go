package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
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

	// create the http client
	client := http.Client{}

	// create the request
	req, err := http.NewRequest("GET", prometheusUrl+"/api/v1/query_range", nil)
	if err != nil {
		panic(err)
	}

	// set the query parameters
	q := req.URL.Query()
	q.Add("query", strings.Join(metrics, ","))
	q.Add("start", from)
	q.Add("end", to)
	q.Add("step", interval)
	req.URL.RawQuery = q.Encode()

	// set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// make the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// check the response status code
	if resp.StatusCode != http.StatusOK {
		panic("Error: " + resp.Status)
	}

	// print the response body
	_, err = os.Stdout.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
}
