package configs

// Default returns a config instance with default values.
func Default() Config {
	return Config{
		RequestTimeout: 30,
		DataDir:        "data",
		PrometheusURL:  "http://localhost:9090",
		Steps:          "5s",
	}
}
