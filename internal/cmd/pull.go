package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/amirhnajafiz/promdigger/internal/client"
	"github.com/amirhnajafiz/promdigger/internal/configs"
	"github.com/amirhnajafiz/promdigger/pkg/files"

	"github.com/spf13/cobra"
)

// PullCMD command gets the records of one single Prometheus query.
type PullCMD struct {
	// public fields
	RootCMD *RootCMD

	// private fields
	query       string
	queryOutput string
	cfg         *configs.Config

	// query fields
	queryStep  time.Duration
	queryStart time.Time
	queryEnd   time.Time
}

// Command builds and returns the cobra command of PullCMD.
func (p *PullCMD) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "pull",
		Short: "Pull a single query records",
		Long:  "Pull records of a single Prometheus query",
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.initVars(); err != nil {
				log.Fatalf("init variables failed: %v\n", err)
			}

			p.main()
		},
	}

	command.
		PersistentFlags().
		StringVar(&p.query, "query", "node_cpu", "prometheus query")
	command.
		PersistentFlags().
		StringVar(&p.queryOutput, "output", "node_cpu", "output file name")

	return command
}

func (p *PullCMD) initVars() error {
	// set the config path
	cfgPath := p.RootCMD.ConfigPath
	if cfgPath == "~/.promdigger/config.json" {
		home, err := os.UserHomeDir() // works across platforms
		if err == nil {
			cfgPath = path.Join(home, ".promdigger", "config.json")
		}
	}

	// initialize the configuration
	cfg, err := configs.LoadConfigs(cfgPath)
	if err != nil {
		return err
	}

	p.cfg = cfg

	// create the output directory
	if err := files.CheckDir(p.cfg.DataDir); err != nil {
		return err
	}

	// convert steps to duration
	du, err := time.ParseDuration(p.cfg.Step)
	if err != nil {
		return fmt.Errorf("invalid duration for step `%s`: %v", p.cfg.Step, err)
	}
	p.queryStep = du

	// convert from and to into time.Time
	startDT, err := time.Parse(time.RFC3339, p.RootCMD.StartFlag)
	if err != nil {
		return fmt.Errorf("invalid time for start `%s`: %v", p.RootCMD.StartFlag, err)
	}
	p.queryStart = startDT.UTC()

	endDT, err := time.Parse(time.RFC3339, p.RootCMD.EndFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", p.RootCMD.EndFlag, err)
	}
	p.queryEnd = endDT.UTC()

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", p.RootCMD.StartFlag, p.RootCMD.EndFlag)
	}

	return nil
}

func (p *PullCMD) main() {
	log.Printf("query `%s` start.\n", p.query)

	// create a client instance
	promClient := client.Client{
		Series:     p.cfg.EstimatedSeriesCount,
		Timeout:    p.cfg.RequestTimeout,
		URL:        p.cfg.PrometheusURL,
		Query:      p.query,
		Step:       p.cfg.Step,
		OutputPath: path.Join(p.cfg.DataDir, p.queryOutput),
	}

	// call split range
	ranges := promClient.TimeRanges(p.queryStart, p.queryEnd, p.queryStep)

	// on ranges call pull and save
	for i := 0; i < len(ranges)-1; i++ {
		start := ranges[i]
		end := ranges[i+1]

		// call pull to fetch metrics with optimized request
		response, err := promClient.Pull(start, end)
		if err != nil {
			log.Fatalf("failed reaching Prometheus: %v\n", err)
		}

		if p.RootCMD.JSONOut {
			// call json export
			if err := promClient.JSONExport(response); err != nil {
				log.Fatalf("failed to export JSON: %v\n", err)
			}
		}
		if p.RootCMD.CSVOut {
			// convert to qqr
			qqr, err := promClient.JSONToQRR(response)
			if err != nil {
				log.Fatalf("failed QQR convert: %v\n", err)
			}

			// extract extra labels
			labels := make([]string, 0)
			if p.cfg.AddExtraCSVLabels {
				labels = append(labels, promClient.ExtractLabels(qqr)...)
			}

			// export to csv
			if err := promClient.CSVExport(qqr, labels...); err != nil {
				log.Fatalf("failed to export CSV: %v\n", err)
			}
		}

		// dump to STDOUT if there is not export flag provided
		if !p.RootCMD.CSVOut && !p.RootCMD.JSONOut {
			fmt.Println(string(response))
		}
	}

	log.Printf("query `%s` completed.\n", p.query)
}
