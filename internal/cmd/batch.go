package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/client"
	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"
	"github.com/amirhnajafiz/prometheus-digger/pkg/models"

	"github.com/spf13/cobra"
)

type BatchCMD struct {
	// public fields
	RootCMD *RootCMD

	// private fields
	inputPath string
	batch     *models.Batch
	cfg       *configs.Config

	// query fields
	queryStep  time.Duration
	queryStart time.Time
	queryEnd   time.Time
}

// Command builds and returns the cobra command of BatchCMD.
func (b *BatchCMD) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "batch",
		Short: "Batch pull multiple queries records",
		Long:  "Batch pull records of Prometheus queries",
		Run: func(cmd *cobra.Command, args []string) {
			if err := b.initVars(); err != nil {
				panic(err)
			}

			b.main()
		},
	}

	command.
		PersistentFlags().
		StringVar(&b.inputPath, "input", "queries.txt", "prometheus queries in a txt file")

	return command
}

func (b *BatchCMD) initVars() error {
	// initialize the configuration
	cfg, err := configs.LoadConfigs(b.RootCMD.ConfigPath)
	if err != nil {
		return err
	}

	b.cfg = cfg

	// create the output directory
	if err := files.CheckDir(b.cfg.DataDir); err != nil {
		return err
	}

	// check the input file
	if ok := files.CheckFile(b.inputPath); !ok {
		return fmt.Errorf("failed to find input file `%s`", b.inputPath)
	}

	// convert steps to duration
	du, err := time.ParseDuration(b.cfg.Step)
	if err != nil {
		return fmt.Errorf("invalid duration for step `%s`: %v", b.cfg.Step, err)
	}
	b.queryStep = du

	// convert from and to into time.Time
	startDT, err := time.Parse(time.RFC3339, b.RootCMD.StartFlag)
	if err != nil {
		return fmt.Errorf("invalid time for start `%s`: %v", b.RootCMD.StartFlag, err)
	}
	b.queryStart = startDT

	endDT, err := time.Parse(time.RFC3339, b.RootCMD.EndFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", b.RootCMD.EndFlag, err)
	}
	b.queryEnd = endDT

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", b.RootCMD.StartFlag, b.RootCMD.EndFlag)
	}

	// read the input and store in batch
	b.batch = models.NewBatch()
	if err := models.FillBatchFromFile(b.inputPath, b.batch); err != nil {
		return err
	}

	return nil
}

func (b *BatchCMD) main() {
	for k, v := range b.batch.Records {
		// create a client instance
		promClient := client.Client{
			Series:     b.cfg.EstimatedSeriesCount,
			Timeout:    b.cfg.RequestTimeout,
			URL:        b.cfg.PrometheusURL,
			Query:      v,
			Step:       b.cfg.Step,
			OutputPath: path.Join(b.cfg.DataDir, k),
		}

		// call split range
		ranges := promClient.TimeRanges(b.queryStart, b.queryEnd, b.queryStep)

		// on ranges call pull and save
		for i := 0; i < len(ranges)-1; i++ {
			start := ranges[i]
			end := ranges[i+1]

			// call pull to fetch metrics with optimized request
			response, err := promClient.Pull(start, end)
			if err != nil {
				panic(err)
			}

			if b.RootCMD.JSONOut {
				// call json export
				if err := promClient.JSONExport(response); err != nil {
					panic(err)
				}
			}
			if b.RootCMD.CSVOut {
				// convert to qqr
				qqr, err := promClient.JSONToQRR(response)
				if err != nil {
					panic(err)
				}

				// extract extra labels
				labels := make([]string, 0)
				if b.cfg.AddExtraCSVLabels {
					labels = append(labels, promClient.ExtractLabels(qqr)...)
				}

				// export to csv
				if err := promClient.CSVExport(qqr, labels...); err != nil {
					panic(err)
				}
			}

			// dump to STDOUT if there is not export flag provided
			if !b.RootCMD.CSVOut && !b.RootCMD.JSONOut {
				fmt.Println(string(response))
			}
		}
	}
}
