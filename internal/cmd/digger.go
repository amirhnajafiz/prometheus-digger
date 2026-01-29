package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/client"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

// Digger is the main handler for fetching the metrics from Prometheus API.
type Digger struct {
	Timeout int
	Metric  string
	Name    string
	URL     string
	Step    time.Duration
	From    time.Time
	To      time.Time
}

// NewDigger returns an instance of digger.
func NewDigger(
	cfg *configs.Config,
	metric,
	name,
	from,
	to string,
) (*Digger, error) {
	// create a new digger instance
	instance := &Digger{
		Timeout: cfg.RequestTimeout,
		Metric:  metric,
		Name:    name,
		URL:     cfg.PrometheusURL,
	}

	// convert steps to duration
	du, err := time.ParseDuration(cfg.Steps)
	if err != nil {
		return nil, fmt.Errorf("invalid duration for steps `%s`: %v", cfg.Steps, err)
	}
	instance.Step = du

	// convert from and to into time.Time
	fd, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid time for from `%s`: %v", from, err)
	}
	instance.From = fd

	td, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid time for to `%s`: %v", to, err)
	}
	instance.To = td

	return instance, nil
}

// Dig collects a metric from Prometheus API.
func (d *Digger) Dig() error {
	var (
		err      error
		response []byte
		tranges  []time.Time
	)

	// get expected datapoints
	if dp := client.GetDataPoints(d.From, d.To, d.Step); dp > 1000 {
		tranges = client.SplitTimeRange(d.From, d.To, d.Step, 1000)
	} else {
		tranges = []time.Time{d.From, d.To}
	}

	// make sure to send long requests as GET
	var handler func(url, metric, step string, from, to time.Time) ([]byte, error)
	if len(d.Metric) < 1024 {
		handler = client.FetchMetricByGET
	} else {
		handler = client.FetchMetricByPOST
	}

	// loop over time ranges and submit the requests
	for i := 0; i < len(tranges)-1; i++ {
		from := tranges[i]
		to := tranges[i+1]

		// call handler function
		response, err = handler(
			d.URL,
			d.Metric,
			d.Step.String(),
			from,
			to,
		)

		// check for errors
		if err != nil {
			return fmt.Errorf("failed to fetch metrics: %v", err)
		}

		// write the output to JSON file
		if err := files.WriteFile(path.Join("out", d.Name), response); err != nil {
			return fmt.Errorf("failed to save metrics: %v", err)
		}
	}

	return nil
}
