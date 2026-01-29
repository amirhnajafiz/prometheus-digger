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
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

// Digger is the main handler for fetching the metrics from Prometheus API.
type Digger struct {
	ESC         int
	HTTPTimeout int
	PromMetric  string
	PromURL     string
	OutputPath  string

	queryDataPoints int
	queryStep       time.Duration
	queryFrom       time.Time
	queryTo         time.Time
	queryRange      []time.Time
}

// NewDigger returns an instance of digger.
func (d *Digger) Validate(from, to, step string) error {
	// convert steps to duration
	du, err := time.ParseDuration(step)
	if err != nil {
		return fmt.Errorf("invalid duration for step `%s`: %v", step, err)
	}
	d.queryStep = du

	// convert from and to into time.Time
	fd, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return fmt.Errorf("invalid time for from `%s`: %v", from, err)
	}
	d.queryFrom = fd

	td, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return fmt.Errorf("invalid time for to `%s`: %v", to, err)
	}
	d.queryTo = td

	// check the from and to range
	if fd.After(td) {
		return fmt.Errorf("to datetime must be after from: %s - %s", from, to)
	}

	// get expected datapoints
	if dp := client.GetDataPoints(d.queryFrom, d.queryTo, d.queryStep, d.ESC); dp > 1000 {
		d.queryRange = client.SplitTimeRange(d.queryFrom, d.queryTo, d.queryStep, 1000)
	} else {
		d.queryRange = []time.Time{d.queryFrom, d.queryTo}
	}

	return nil
}

// Dig collects a metric from Prometheus API.
func (d *Digger) Dig() error {
	// make sure to send long requests as GET
	var handler func(url, metric, step string, from, to time.Time) ([]byte, error)
	if len(d.PromMetric) < 1024 {
		handler = client.FetchMetricByGET
	} else {
		handler = client.FetchMetricByPOST
	}

	// loop over time ranges and submit the requests
	for i := 0; i < len(d.queryRange)-1; i++ {
		from := d.queryRange[i]
		to := d.queryRange[i+1]

		// call handler function
		response, err := handler(
			d.PromURL,
			d.PromMetric,
			d.queryStep.String(),
			from,
			to,
		)

		// check for errors
		if err != nil {
			return fmt.Errorf("failed to fetch metrics: %v", err)
		}

		// write the output to JSON file
		if err := d.writeQueryRangeCSV(response, path.Join(d.OutputPath+".csv")); err != nil {
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

	var (
		file       *os.File
		err        error
		appendMode bool
	)

	// check if file exists
	if files.CheckFile(outputPath) {
		file, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		appendMode = true
	} else {
		file, err = os.Create(outputPath)
		appendMode = false
	}

	// check errors
	if err != nil {
		return err
	}
	defer file.Close()

	// create a new csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// headers to be added only if file created for the first time
	if !appendMode {
		header := append([]string{"timestamp", "metric", "value"}, labels...)
		if err := writer.Write(header); err != nil {
			return err
		}
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
