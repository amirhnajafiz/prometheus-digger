package client

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/datapoints"
)

const dbLimit = 1000

// Client is responsible for managing requests to the Prometheus API.
type Client struct {
	Series  int
	Timeout int
	URL     string
	Query   string
	Step    string
}

// TimeRanges split the query from `start` to `end`, making sure that each
// request contains less that `dbLimit` datapoints in its response.
func (c *Client) TimeRanges(start, end time.Time, step time.Duration) []time.Time {
	// get expected datapoints
	var tranges []time.Time
	if dp := datapoints.GetDataPoints(start, end, step, c.Series); dp > dbLimit {
		tranges = datapoints.SplitTimeRange(start, end, step, dbLimit, dp)
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
