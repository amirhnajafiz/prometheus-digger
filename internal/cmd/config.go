package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/spf13/cobra"
)

// ConfigCMD command checks the current config of the app.
type ConfigCMD struct {
	// public fields
	RootCMD *RootCMD

	// private fields
	cfg *configs.Config
}

// Command builds and returns the cobra command of ConfigCMD.
func (c *ConfigCMD) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "View config",
		Long:  "View the current configuration values",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.initVars(); err != nil {
				log.Fatalf("init variables failed: %v\n", err)
			}

			c.main()
		},
	}

	return command
}

func (c *ConfigCMD) initVars() error {
	// set the config path
	cfgPath := c.RootCMD.ConfigPath
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

	c.cfg = cfg

	return nil
}

func (c *ConfigCMD) main() {
	fmt.Printf(
		"Values:\n\tAdding Extra CSV Labels: %v\n\tData Dir: %s\n\tEstimated Series Count: %d\n\tPrometheus URL: %s\n\tHTTP Timeout: %d\n\tStep: %s\n",
		c.cfg.AddExtraCSVLabels,
		c.cfg.DataDir,
		c.cfg.EstimatedSeriesCount,
		c.cfg.PrometheusURL,
		c.cfg.RequestTimeout,
		c.cfg.Step,
	)
}
