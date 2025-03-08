package internal

// Worker fetches metrics from Prometheus.
// It takes the Prometheus URL, metric name, start time, end time, and interval as parameters.
// It returns the response body as a string and an error if any.
// It uses the newHttpGetRequest and fetchMetrics functions from the network package.
// It sets the query parameters for the request and the headers for the request.
// It returns the response body as a string and an error if any.
func Worker(url, metric, from, to, interval string) (string, error) {
	// create HTTP GET request
	req, err := newHttpGetRequest(url)
	if err != nil {
		return "", err
	}

	// set query parameters
	q := req.URL.Query()
	q.Add("query", metric)
	q.Add("start", from)
	q.Add("end", to)
	q.Add("step", interval)
	req.URL.RawQuery = q.Encode()

	// set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// fetch metrics
	rsp, err := fetchMetrics(req)
	if err != nil {
		return "", err
	}

	return rsp, nil
}
