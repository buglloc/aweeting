package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:          "start",
	SilenceUsage: true,
	Short:        "Start API srv",
	RunE: func(cmd *cobra.Command, args []string) error {
		runtime, err := cfg.NewRuntime()
		if err != nil {
			return fmt.Errorf("create runtime: %w", err)
		}

		updater, err := runtime.NewAwtrixUpdater()
		if err != nil {
			return fmt.Errorf("create updater: %w", err)
		}

		tick, err := runtime.NewTicker()
		if err != nil {
			return fmt.Errorf("create ticker: %w", err)
		}

		errChan := make(chan error, 1)
		okChan := make(chan struct{}, 1)
		go func() {
			if err := tick.Start(updater.Update); err != nil {
				errChan <- fmt.Errorf("failed to start application: %w", err)
			} else {
				okChan <- struct{}{}
			}
		}()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		defer log.Info().Msg("stopped")

		select {
		case <-sigChan:
			log.Info().Msg("shutting down gracefully by signal")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			tick.Stop(ctx)
			return nil
		case <-okChan:
			return nil
		case err := <-errChan:
			return err
		}
	},
}
