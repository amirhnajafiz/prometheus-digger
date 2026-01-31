package cmd

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/amirhnajafiz/prometheus-digger/pkg/files"

	"github.com/spf13/cobra"
)

type PullCMD struct {
	// public fields
	ConfigPath string

	// private fields
	cfg       *configs.Config
	startFlag string
	endFlag   string

	// query fields
	queryStep  time.Duration
	queryStart time.Time
	queryEnd   time.Time
}

func (p *PullCMD) Command() *cobra.Command {
	return &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.initVars(); err != nil {
				panic(err)
			}

			p.main()
		},
	}
}

func (p *PullCMD) initVars() error {
	// initialize the configuration
	cfg, err := configs.LoadConfigs(p.ConfigPath)
	if err != nil {
		return err
	}

	// create the output directory
	if err := files.CheckDir(cfg.DataDir); err != nil {
		return err
	}

	// convert steps to duration
	du, err := time.ParseDuration(p.cfg.Step)
	if err != nil {
		return fmt.Errorf("invalid duration for step `%s`: %v", p.cfg.Step, err)
	}
	p.queryStep = du

	// convert from and to into time.Time
	startDT, err := time.Parse(time.RFC3339, p.startFlag)
	if err != nil {
		return fmt.Errorf("invalid time for start `%s`: %v", p.startFlag, err)
	}
	p.queryStart = startDT

	endDT, err := time.Parse(time.RFC3339, p.endFlag)
	if err != nil {
		return fmt.Errorf("invalid time for end `%s`: %v", p.endFlag, err)
	}
	p.queryEnd = endDT

	// check the from and to range
	if startDT.After(endDT) {
		return fmt.Errorf("to datetime must be after from: %s - %s", p.startFlag, p.endFlag)
	}

	return nil
}

func (p *PullCMD) main() {}
