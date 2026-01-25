package client

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/pkg/networking"
)

// calling the query range API on the Prometheus
const API = "/api/v1/query_range"

func FetchMetricByGET(url, metric, step string, from, to time.Time) ([]byte, error) {
	// create HTTP GET request
	req, err := networking.NewHttpGetRequest(url + API)
	if err != nil {
		return nil, fmt.Errorf("build GET request failed: %v", err)
	}

	// set query parameters
	q := req.URL.Query()
	q.Add("query", metric)
	q.Add("start", from.Format(time.RFC3339))
	q.Add("end", to.Format(time.RFC3339))
	q.Add("step", step)
	req.URL.RawQuery = q.Encode()

	// set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// fetch metrics by sending the http request
	return networking.HTTPSend(req)
}

func FetchMetricByPOST(url, metric, step string, from, to time.Time) ([]byte, error) {
	// form the request body
	body := fmt.Sprintf(
		"query=%s&start=%s&end=%s&step=%s",
		metric,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
		step,
	)

	// create HTTP POST request
	req, err := networking.NewHttpPostRequest(url+API, []byte(body))
	if err != nil {
		return nil, fmt.Errorf("build POST request failed: %v", err)
	}

	// set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// fetch metrics by sending the http request
	return networking.HTTPSend(req)
}
