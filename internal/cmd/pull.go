package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/client"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"

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
				panic(err)
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
	// initialize the configuration
	cfg, err := configs.LoadConfigs(p.RootCMD.ConfigPath)
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
	p.queryStart = startDT

	endDT, err := time.Parse(time.RFC3339, p.RootCMD.EndFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", p.RootCMD.EndFlag, err)
	}
	p.queryEnd = endDT

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", p.RootCMD.StartFlag, p.RootCMD.EndFlag)
	}

	return nil
}

func (p *PullCMD) main() {
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
			panic(err)
		}

		if p.RootCMD.JSONOut {
			if err := promClient.JSONExport(response); err != nil {
				panic(err)
			}
		} else if p.RootCMD.CSVOut {
			// convert to qqr
			qqr, err := promClient.JSONToQRR(response)
			if err != nil {
				panic(err)
			}

			// extract extra labels
			labels := make([]string, 0)
			if p.cfg.AddExtraCSVLabels {
				labels = append(labels, promClient.ExtractLabels(qqr)...)
			}

			// export to csv
			if err := promClient.CSVExport(qqr, labels...); err != nil {
				panic(err)
			}
		} else {
			fmt.Println(string(response))
		}
	}
}
