package internal

import (
	"io"
	"net/http"
)

// newHttpGetRequest creates a new HTTP GET request with the given URL.
func newHttpGetRequest(url string) (*http.Request, error) {
	return http.NewRequest("GET", url, nil)
}

// fetchMetrics sends the given HTTP request and returns the response body as a string.
func fetchMetrics(req *http.Request) (string, error) {
	// create the http client
	client := http.Client{}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
