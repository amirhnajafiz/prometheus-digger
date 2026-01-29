package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/client"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
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
		if err := d.writeQueryRangeCSV(response, path.Join("out", d.Name+".csv")); err != nil {
			return fmt.Errorf("failed to save metrics: %v", err)
		}
	}

	return nil
}

// parses Prometheus query_range JSON bytes and writes CSV.
func (d *Digger) writeQueryRangeCSV(apiBytes []byte, outputPath string) error {
	var resp client.QueryRangeResponse
	if err := json.Unmarshal(apiBytes, &resp); err != nil {
		return fmt.Errorf("invalid prometheus response: %w", err)
	}

	if resp.Status != "success" || resp.Data.ResultType != "matrix" {
		return fmt.Errorf("unexpected prometheus response")
	}

	// collect all label keys (for CSV header)
	labelSet := map[string]struct{}{}
	for _, series := range resp.Data.Result {
		for k := range series.Metric {
			labelSet[k] = struct{}{}
		}
	}

	labels := make([]string, 0, len(labelSet))
	for k := range labelSet {
		labels = append(labels, k)
	}
	sort.Strings(labels)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header
	header := append([]string{"timestamp", "metric", "value"}, labels...)
	if err := writer.Write(header); err != nil {
		return err
	}

	// rows
	for _, series := range resp.Data.Result {
		metricName := series.Metric["__name__"]

		for _, v := range series.Values {
			ts := fmt.Sprintf("%.0f", v[0].(float64))
			value := v[1].(string)

			row := []string{ts, metricName, value}
			for _, label := range labels {
				row = append(row, series.Metric[label])
			}

			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}
