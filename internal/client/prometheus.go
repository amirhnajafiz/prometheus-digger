package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// calling the query range API on the Prometheus
const API = "/api/v1/query_range"

func FetchMetricByGET(url, metric, step string, from, to time.Time) ([]byte, error) {
	// create HTTP GET request
	req, err := http.NewRequest("GET", url+API, nil)
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
	return sendHTTPReqeust(req)
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
	req, err := http.NewRequest("POST", url+API, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("build POST request failed: %v", err)
	}

	// set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// fetch metrics by sending the http request
	return sendHTTPReqeust(req)
}

// sends the given HTTP request and returns the response body.
func sendHTTPReqeust(req *http.Request) ([]byte, error) {
	// create the http client
	client := http.Client{}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed HTTP: [`%d`] %s", resp.StatusCode, string(body))
	}

	return body, nil
}
