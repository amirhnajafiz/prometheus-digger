package client

import (
	"fmt"
	"time"
)

// Client is responsible for managing requests to the Prometheus API.
type Client struct {
	URL     string
	Query   string
	Step    string
	Timeout int
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
