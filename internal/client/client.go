package client

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/datapoints"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

const (
	DatapointsLimit = 1000
	API             = "/api/v1/query_range"
)

// Client is responsible for managing requests to the Prometheus API.
type Client struct {
	Series     int
	Timeout    int
	URL        string
	Query      string
	Step       string
	OutputPath string
}

// TimeRanges split the query from `start` to `end`, making sure that each
// request contains less that `dbLimit` datapoints in its response.
func (c *Client) TimeRanges(start, end time.Time, step time.Duration) []time.Time {
	// get expected datapoints
	var tranges []time.Time
	if dp := datapoints.GetDataPoints(start, end, step, c.Series); dp > DatapointsLimit {
		tranges = datapoints.SplitTimeRange(start, end, step, DatapointsLimit, dp)
	} else {
		tranges = []time.Time{start, end}
	}

	return tranges
}

// Pull the records of the given query in from `start` to `end`.
func (c *Client) Pull(start, end time.Time) ([]byte, error) {
	// make sure to send long requests as POST and short requests as GET
	var handler func(time.Time, time.Time) ([]byte, error)
	if len(c.Query) < 1024 {
		handler = c.GET
	} else {
		handler = c.POST
	}

	// call handler function
	response, err := handler(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %v", err)
	}

	return response, nil
}

// JSONToQRR, accepts a JSON response from API and converts it to a
// QueryRangeResponse object.
func (c *Client) JSONToQRR(bytes []byte) (*QueryRangeResponse, error) {
	var resp QueryRangeResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return nil, fmt.Errorf("invalid prometheus response: %w", err)
	}

	if resp.Status != "success" || resp.Data.ResultType != "matrix" {
		return nil, fmt.Errorf("unexpected prometheus response")
	}

	return &resp, nil
}

// ExtractLabels, gets a QueryRangeResponse and returns a list of labels as str.
func (c *Client) ExtractLabels(qrr *QueryRangeResponse) []string {
	var labels []string

	labelSet := map[string]struct{}{}
	for _, series := range qrr.Data.Result {
		for k := range series.Metric {
			labelSet[k] = struct{}{}
		}
	}

	labels = make([]string, 0, len(labelSet))
	for k := range labelSet {
		labels = append(labels, k)
	}

	sort.Strings(labels)

	return labels
}

// JSONExport, writes the JSON response from Prometheus API into a file.
func (c *Client) JSONExport(response []byte) error {
	return files.WriteFile(c.OutputPath+".json", response)
}

// CSVExport, parses Prometheus query_range JSON bytes and writes them into
// a CSV file with given output path.
func (c *Client) CSVExport(qrr *QueryRangeResponse, labels ...string) error {
	var (
		file       *os.File
		err        error
		appendMode bool
	)

	// check if file exists
	path := c.OutputPath + ".csv"
	if appendMode = files.CheckFile(path); appendMode {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	} else {
		file, err = os.Create(path)
	}

	// check errors
	if err != nil {
		return fmt.Errorf("failed to open output file `%s`: %v", c.OutputPath, err)
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

	// insert rows
	for _, series := range qrr.Data.Result {
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
