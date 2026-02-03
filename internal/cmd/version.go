package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const appVersion = "v0.1.0"

// VCMD is just a struct for reading the version of application.
type VCMD struct{}

// Command returns the cobra command of VCMD.
func (v *VCMD) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "app version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(appVersion)
		},
	}
}
