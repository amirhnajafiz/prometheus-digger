package cmd

import (
	"github.com/spf13/cobra"
)

// RootCMD is the root cobra command handler.
type RootCMD struct {
	ConfigPath string
	StartFlag  string
	EndFlag    string
}

// Command builds and returns the cobra command of RootCMD.
func (r *RootCMD) Command() *cobra.Command {
	command := &cobra.Command{}

	command.
		PersistentFlags().
		StringVarP(&r.ConfigPath, "config-path", "c", "config.json", "path to a JSON configuration file")
	command.
		PersistentFlags().
		StringVarP(&r.StartFlag, "start", "s", "2026-01-31T18:33:44-05:00", "query start time")
	command.
		PersistentFlags().
		StringVarP(&r.EndFlag, "end", "e", "2026-01-31T18:34:44-05:00", "query end time")

	return command
}
