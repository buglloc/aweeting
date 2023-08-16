package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootArgs struct {
	Verbose bool
	Source  string
}

var rootCmd = &cobra.Command{
	Use:          "aweeting",
	SilenceUsage: true,
	Short:        "Meeting timer app for awtrix-light",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.BoolVar(&rootArgs.Verbose, "verbose", os.Getenv("AW_VERBOSE") != "", "verbose")
	flags.StringVar(&rootArgs.Source, "source", os.Getenv("AW_SOURCE"), "source")

	rootCmd.AddCommand(
		eventsCmd,
		startCmd,
	)
}
