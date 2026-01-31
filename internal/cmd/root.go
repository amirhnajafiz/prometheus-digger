package cmd

import "github.com/spf13/cobra"

// RootCMD is the root cobra command handler.
type RootCMD struct {
	ConfigPath string
}

// Command builds and returns the cobra command of RootCMD.
func (r *RootCMD) Command() *cobra.Command {
	command := &cobra.Command{}

	command.
		PersistentFlags().
		StringVarP(&r.ConfigPath, "config-path", "cp", "config.json", "path to a JSON configuration file")

	return command
}
