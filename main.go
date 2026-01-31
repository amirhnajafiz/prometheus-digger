package main

import "github.com/amirhnajafiz/prometheus-digger/internal/cmd"

func main() {
	// create root command
	rcm := cmd.RootCMD{}
	root := rcm.Command()

	// add sub-commands
	root.AddCommand(
		(&cmd.PullCMD{
			ConfigPath: rcm.ConfigPath,
			StartFlag:  rcm.StartFlag,
			EndFlag:    rcm.EndFlag,
		}).Command(),
	)

	// execute the command
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
