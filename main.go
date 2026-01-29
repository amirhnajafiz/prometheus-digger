package main

import (
	"flag"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/cmd"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
)

func main() {
	// register the flags
	var (
		FlagMetric         = flag.String("metric", "node_cpu_seconds_total", "Metric to fetch from Prometheus")
		FlagName           = flag.String("name", "node_cpu", "Output name")
		FlagTimeFrom       = flag.String("from", time.Now().Format(time.RFC3339), "The start timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagTimeTo         = flag.String("to", time.Now().Add(5*time.Minute).Format(time.RFC3339), "The end timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagConfigFilePath = flag.String("config", "config.json", "Path to the configuration file")
	)

	flag.Parse()

	// initialize the configuration
	cfg, err := configs.LoadConfigs(*FlagConfigFilePath)
	if err != nil {
		panic(err)
	}

	// create a digger instance
	digger, err := cmd.NewDigger(
		cfg,
		*FlagMetric,
		*FlagName,
		*FlagTimeFrom,
		*FlagTimeTo,
	)
	if err != nil {
		panic(err)
	}

	// call the digger
	if err := digger.Dig(); err != nil {
		panic(err)
	}
}
