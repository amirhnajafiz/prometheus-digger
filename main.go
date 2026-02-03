package main

import "github.com/amirhnajafiz/prometheus-digger/internal/cmd"

func main() {
	// create root command
	rcm := cmd.RootCMD{}
	root := rcm.Command()

	// add sub-commands
	root.AddCommand(
		(&cmd.ConfigCMD{
			RootCMD: &rcm,
		}).Command(),
		(&cmd.PullCMD{
			RootCMD: &rcm,
		}).Command(),
		(&cmd.BatchCMD{
			RootCMD: &rcm,
		}).Command(),
	)

	// execute the command
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
