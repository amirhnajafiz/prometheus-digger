package pkg

import (
	"fmt"
	"io"
	"net/http"
)

// NewHttpGetRequest creates a new HTTP GET request with the given URL.
func NewHttpGetRequest(url string) (*http.Request, error) {
	return http.NewRequest("GET", url, nil)
}

// FetchMetrics sends the given HTTP request and returns the response body as a string.
func FetchMetrics(req *http.Request) ([]byte, error) {
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
		return nil, fmt.Errorf("prometheus API response code: %d, %s", resp.StatusCode, string(body))
	}

	return body, nil
}
