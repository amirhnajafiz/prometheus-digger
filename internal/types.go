package internal

// config is a module that holds the configuration of the application.
type config struct {
	URL      string  `json:"url"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	PoolSize int     `json:"pool_size"`
	Queries  []query `json:"queries"`
}

// query is a module that holds the query information.
type query struct {
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Interval string `json:"interval"`
}
