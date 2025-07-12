package models

// Query is a module that holds the query information.
type Query struct {
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Interval string `json:"interval"`
}
