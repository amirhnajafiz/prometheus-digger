package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/amirhnajafiz/prometheus-digger/internal/configs"
	"github.com/spf13/cobra"
)

// HealthCMD command checks the status of Prometheus server.
type HealthCMD struct {
	// public fields
	RootCMD *RootCMD

	// private fields
	cfg *configs.Config
}

// Command builds and returns the cobra command of ConfigCMD.
func (h *HealthCMD) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "health",
		Short: "Check Prometheus health",
		Long:  "Check Prometheus server reachability and status",
		Run: func(cmd *cobra.Command, args []string) {
			if err := h.initVars(); err != nil {
				panic(err)
			}

			h.main()
		},
	}

	return command
}

func (h *HealthCMD) initVars() error {
	// set the config path
	cfgPath := h.RootCMD.ConfigPath
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

	h.cfg = cfg

	return nil
}

func (h *HealthCMD) main() {
	log.Printf("checking %s ...\n", h.cfg.PrometheusURL)

	hurl := fmt.Sprintf("%s/healthy", h.cfg.PrometheusURL)
	if resp, err := http.Get(hurl); err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(errors.New("Prometheus returnes not OK!"))
	}

	rurl := fmt.Sprintf("%s/ready", h.cfg.PrometheusURL)
	if resp, err := http.Get(rurl); err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(errors.New("Prometheus returnes not Ready!"))
	}

	log.Println("Prometheus is OK and Ready.")
}
