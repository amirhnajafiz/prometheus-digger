package main

import (
	"flag"

	"github.com/amirhnajafiz/prometheus-digger/internal/cmd"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
)

func main() {
	// register the flags
	var (
		FlagMetric         = flag.String("metric", "", "Metric to fetch from Prometheus")
		FlagTimeFrom       = flag.String("from", "", "The start timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagTimeTo         = flag.String("to", "", "The end timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagConfigFilePath = flag.String("config", "config.json", "Path to the configuration file")
	)

	flag.Parse()

	// initialize the configuration
	cfg, err := configs.LoadConfigs(*FlagConfigFilePath)
	if err != nil {
		panic(err)
	}

	// create a digger instance
	digger, err := cmd.NewDigger(cfg, *FlagMetric, *FlagTimeFrom, *FlagTimeTo)
	if err != nil {
		panic(err)
	}

	// call the digger
	if err := digger.Dig(); err != nil {
		panic(err)
	}
}
