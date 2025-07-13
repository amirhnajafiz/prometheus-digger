package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/internal/parser"
	"github.com/amirhnajafiz/prometheus-digger/internal/worker"
)

func help() {
	fmt.Println("Usage: prometheus-digger [options]")
	fmt.Println("Options:")
	fmt.Println("  -config string   Path to the configuration file (default \"config.json\")")
	fmt.Println("  -help            Show this help message")
}

func initCfg(path string) *configs.Config {
	// load configs
	cfg, err := configs.LoadConfigs(path)
	if err != nil {
		log.Fatalf("[ERR] load configs failed: %v\n", err)
	}

	// parse the configuration dates
	var startTime time.Time
	if cfg.From == "" {
		startTime = time.Now()
		cfg.From = parser.ConvertToString(startTime)
	} else {
		startTime, _ = parser.ConvertToTime(cfg.From)
	}

	cfg.To = parser.ConvertToString(parser.ConvertSliceToTime(startTime, cfg.To))
	cfg.From, cfg.To = parser.SortDates(cfg.From, cfg.To)

	return cfg
}

func main() {
	// register the flags
	var (
		configFilePath = flag.String("config", "config.json", "Path to the configuration file")
		helpFlag       = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	// check if help is requested
	if *helpFlag {
		help()
		return
	}

	// initialize the configuration
	cfg := initCfg(*configFilePath)

	// print the configuration
	fmt.Printf("Configuration loaded:\n\tMetrics From: %s\n\tTo: %s\n", cfg.From, cfg.To)
	fmt.Printf("\tTarget: %s\n", cfg.URL)
	fmt.Printf("\tNumber of queries: %d\n\n", len(cfg.Queries))

	// create worker pool
	pool := worker.NewWorkerPool(cfg)
	if pool == nil {
		log.Fatalln("[ERR] failed to create worker pool!")
		return
	}

	// collect metrics
	pool.Collect()

	// wait for all workers
	pool.StopAndWait()

	// print the summary
	pool.PrintStats()
}
