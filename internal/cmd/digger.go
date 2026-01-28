package cmd

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/client"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

// Digger is the main handler for fetching the metrics from Prometheus API.
type Digger struct {
	Timeout int
	Metric  string
	URL     string
	Steps   time.Duration
	From    time.Time
	To      time.Time
}

// NewDigger returns an instance of digger.
func NewDigger(
	cfg *configs.Config,
	metric,
	from,
	to string,
) (*Digger, error) {
	// create a new digger instance
	instance := &Digger{
		Timeout: cfg.RequestTimeout,
		Metric:  metric,
		URL:     cfg.PrometheusURL,
	}

	// convert steps to duration
	du, err := time.ParseDuration(cfg.Steps)
	if err != nil {
		return nil, fmt.Errorf("invalid duration for steps `%s`: %v", cfg.Steps, err)
	}
	instance.Steps = du

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
	)

	// make sure to send long requests as GET
	if len(d.Metric) < 1024 {
		response, err = client.FetchMetricByGET(
			d.URL,
			d.Metric,
			d.Steps.String(),
			d.From,
			d.To,
		)
	} else {
		response, err = client.FetchMetricByPOST(
			d.URL,
			d.Metric,
			d.Steps.String(),
			d.From,
			d.To,
		)
	}

	// check for errors
	if err != nil {
		return fmt.Errorf("failed to fetch metrics: %v", err)
	}

	// write the output to JSON file
	if err := files.WriteFile("out", response); err != nil {
		return fmt.Errorf("failed to save metrics: %v", err)
	}

	return nil
}
