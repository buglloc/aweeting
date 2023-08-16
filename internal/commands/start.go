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

	"github.com/buglloc/aweeting/internal/awtrix"
	"github.com/buglloc/aweeting/internal/calendar"
	"github.com/buglloc/aweeting/internal/ticker"
)

var startArgs struct {
	Upstream string
	Topic    string
	Username string
	Password string
}

var startCmd = &cobra.Command{
	Use:          "start",
	SilenceUsage: true,
	Short:        "Start API srv",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := calendar.NewICal(rootArgs.Source)
		if err != nil {
			return fmt.Errorf("create calendar: %w", err)
		}

		updater, err := awtrix.NewMqttUpdater(awtrix.MqttConfig{
			Upstream: startArgs.Upstream,
			Topic:    startArgs.Topic,
			Username: startArgs.Username,
			Password: startArgs.Password,
			Icons: awtrix.Set{
				Zero:     "11899",
				Upcoming: "11899",
				OnAir:    "11899",
			},
			Colors: awtrix.Set{
				Zero:     "#ffffff",
				Upcoming: "#ffffff",
				OnAir:    "#e60000",
			},
		})
		if err != nil {
			return fmt.Errorf("create updater: %w", err)
		}

		tick := ticker.NewConstTicker(cl)

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

func init() {
	flags := startCmd.PersistentFlags()
	flags.StringVar(&startArgs.Upstream, "mqtt.upstream", os.Getenv("AW_MQTT_UPSTREAM"), "mqtt upstream")
	flags.StringVar(&startArgs.Username, "mqtt.username", os.Getenv("AW_MQTT_USERNAME"), "mqtt username")
	flags.StringVar(&startArgs.Password, "mqtt.password", os.Getenv("AW_MQTT_PASSWORD"), "mqtt password")
	flags.StringVar(&startArgs.Topic, "mqtt.topic", os.Getenv("AW_MQTT_TOPIC"), "mqtt topic")
}
