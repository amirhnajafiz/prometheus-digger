package internal

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
