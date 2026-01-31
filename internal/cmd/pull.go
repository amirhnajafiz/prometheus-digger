package cmd

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"

	"github.com/spf13/cobra"
)

// PullCMD command gets the records of one single Prometheus query.
type PullCMD struct {
	// public fields
	ConfigPath string
	StartFlag  string
	EndFlag    string

	// private fields
	queryFlag string
	cfg       *configs.Config

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
		StringVarP(&p.queryFlag, "query", "q", "node_uptime", "prometheus query")

	return command
}

func (p *PullCMD) initVars() error {
	// initialize the configuration
	cfg, err := configs.LoadConfigs(p.ConfigPath)
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
	startDT, err := time.Parse(time.RFC3339, p.StartFlag)
	if err != nil {
		return fmt.Errorf("invalid time for start `%s`: %v", p.StartFlag, err)
	}
	p.queryStart = startDT

	endDT, err := time.Parse(time.RFC3339, p.EndFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", p.EndFlag, err)
	}
	p.queryEnd = endDT

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", p.StartFlag, p.EndFlag)
	}

	return nil
}

func (p *PullCMD) main() {

}
