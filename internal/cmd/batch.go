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
	"github.com/amirhnajafiz/promdigger/pkg/models"

	"github.com/spf13/cobra"
)

// BatchCMD reads the queries from a file and pulls them sequentially.
type BatchCMD struct {
	// public fields
	RootCMD *RootCMD

	// private fields
	inputPath string
	batch     *models.Batch
	cfg       *configs.Config
	cpool     *client.ClientObjectPool

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
				log.Fatalf("init variables failed: %v\n", err)
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
	// set the config path
	cfgPath := b.RootCMD.ConfigPath
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
	b.queryStart = startDT.UTC()

	endDT, err := time.Parse(time.RFC3339, b.RootCMD.EndFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", b.RootCMD.EndFlag, err)
	}
	b.queryEnd = endDT.UTC()

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", b.RootCMD.StartFlag, b.RootCMD.EndFlag)
	}

	// read the input and store in batch
	b.batch = models.NewBatch()
	if err := models.FillBatchFromFile(b.inputPath, b.batch); err != nil {
		return err
	}

	// create a new client object pool
	b.cpool = client.NewObjectPool()

	return nil
}

func (b *BatchCMD) main() {
	for k, v := range b.batch.Records {
		log.Printf("query `%s` start.\n", v)

		// create a client instance
		promClient := b.cpool.GetClientObj()

		promClient.Series = b.cfg.EstimatedSeriesCount
		promClient.Timeout = b.cfg.RequestTimeout
		promClient.URL = b.cfg.PrometheusURL
		promClient.Query = v
		promClient.Step = b.cfg.Step
		promClient.OutputPath = path.Join(b.cfg.DataDir, k)

		// call split range
		ranges := promClient.TimeRanges(b.queryStart, b.queryEnd, b.queryStep)

		// on ranges call pull and save
		for i := 0; i < len(ranges)-1; i++ {
			start := ranges[i]
			end := ranges[i+1]

			// call pull to fetch metrics with optimized request
			response, err := promClient.Pull(start, end)
			if err != nil {
				log.Fatalf("failed reaching Prometheus: %v\n", err)
			}

			if b.RootCMD.JSONOut {
				// call json export
				if err := promClient.JSONExport(response); err != nil {
					log.Fatalf("failed to export JSON: %v\n", err)
				}
			}
			if b.RootCMD.CSVOut {
				// convert to qqr
				qqr, err := promClient.JSONToQRR(response)
				if err != nil {
					log.Fatalf("failed QQR convert: %v\n", err)
				}

				// extract extra labels
				labels := make([]string, 0)
				if b.cfg.AddExtraCSVLabels {
					labels = append(labels, promClient.ExtractLabels(qqr)...)
				}

				// export to csv
				if err := promClient.CSVExport(qqr, labels...); err != nil {
					log.Fatalf("failed to export CSV: %v\n", err)
				}
			}

			// dump to STDOUT if there is not export flag provided
			if !b.RootCMD.CSVOut && !b.RootCMD.JSONOut {
				fmt.Println(string(response))
			}
		}

		log.Printf("query `%s` completed.\n", v)

		b.cpool.PutClientObj(promClient)
	}
}
