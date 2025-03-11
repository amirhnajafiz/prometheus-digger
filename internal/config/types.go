package config

// Config is a module that holds the configuration of the application.
type Config struct {
	URL      string  `json:"url"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	PoolSize int     `json:"pool_size"`
	Queries  []Query `json:"queries"`
}

// Query is a module that holds the query information.
type Query struct {
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Interval string `json:"interval"`
}
