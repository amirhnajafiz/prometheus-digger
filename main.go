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

	// load configs
	cfg, err := configs.LoadConfigs(*configFilePath)
	if err != nil {
		log.Printf("[ERR] load configs failed: %v\n", err)
		return
	}

	// parse the configuration dates
	startTime := time.Now()
	cfg.From = parser.ConvertToString(parser.ConvertSliceToTime(startTime, cfg.From))
	cfg.To = parser.ConvertToString(parser.ConvertSliceToTime(startTime, cfg.To))

	// print the configuration
	fmt.Printf("Configuration loaded:\nFrom: %s\nTo: %s\n", cfg.From, cfg.To)
	fmt.Printf("Target: %s\n", cfg.URL)

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
