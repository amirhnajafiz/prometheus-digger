package configs

// Default returns a config instance with default values.
func Default() Config {
	return Config{
		AddExtraCSVLabels:    true,
		EstimatedSeriesCount: 1,
		RequestTimeout:       30,
		DataDir:              "data",
		PrometheusURL:        "http://localhost:9090",
		Step:                 "5s",
	}
}
