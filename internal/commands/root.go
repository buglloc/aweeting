package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/buglloc/aweeting/internal/config"
)

var cfg *config.Config

var rootArgs struct {
	Configs []string
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
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	flags := rootCmd.PersistentFlags()
	flags.StringSliceVar(&rootArgs.Configs, "config", nil, "config file")

	rootCmd.AddCommand(
		startCmd,
		eventsCmd,
	)
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig(rootArgs.Configs...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to load config: %v\n", err)
		os.Exit(1)
	}
}

func initLogger() {
	log.SetOutput(os.Stderr)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
