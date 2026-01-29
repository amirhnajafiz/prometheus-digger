package main

import (
	"flag"
	"path"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/cmd"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
)

func main() {
	// register the flags
	var (
		FlagMetric         = flag.String("metric", "node_cpu_seconds_total", "Metric(or Query) to fetch from the Prometheus API")
		FlagName           = flag.String("name", "node_cpu", "Output file name (node_cpu)")
		FlagTimeFrom       = flag.String("from", time.Now().Format(time.RFC3339), "Query start timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagTimeTo         = flag.String("to", time.Now().Add(5*time.Minute).Format(time.RFC3339), "Query end timestamp in RFC3339 (2006-01-02T15:04:05Z07:00)")
		FlagConfigFilePath = flag.String("config", "config.json", "Path to the configuration file")
	)

	flag.Parse()

	// initialize the configuration
	cfg, err := configs.LoadConfigs(*FlagConfigFilePath)
	if err != nil {
		panic(err)
	}

	// create the output directory
	if err := files.CheckDir(cfg.DataDir); err != nil {
		panic(err)
	}

	// create a digger instance and validate the inputs
	digger := cmd.Digger{
		HTTPTimeout: cfg.RequestTimeout,
		PromMetric:  *FlagMetric,
		PromURL:     cfg.PrometheusURL,
		OutputPath:  path.Join(cfg.DataDir, *FlagName),
	}
	if err := digger.Validate(*FlagTimeFrom, *FlagTimeTo, cfg.Steps); err != nil {
		panic(err)
	}

	// call the digger
	if err := digger.Dig(); err != nil {
		panic(err)
	}
}
