package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/config"
	"github.com/amirhnajafiz/prometheus-digger/internal/parser"
	"github.com/amirhnajafiz/prometheus-digger/internal/worker"
)

func main() {
	// register the flags
	var (
		configFilePath = flag.String("config", "config.json", "Path to the configuration file")
	)

	// load configs
	cfg, err := config.LoadConfigs(*configFilePath)
	if err != nil {
		log.Printf("[ERR] load configs failed: %v\n", err)
		return
	}

	// parse the configuration dates
	startTime := time.Now()
	cfg.From = parser.ConvertToString(parser.ConvertSliceToTime(startTime, cfg.From))
	cfg.To = parser.ConvertToString(parser.ConvertSliceToTime(startTime, cfg.To))

	// create worker pool
	pool := worker.NewWorkerPool(cfg)
	if pool == nil {
		fmt.Println("[ERR] failed to create worker pool!")
		return
	}

	// collect metrics
	pool.Collect()

	// wait for all workers
	pool.StopAndWait()
}
