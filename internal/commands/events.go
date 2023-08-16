package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/buglloc/aweeting/internal/calendar"
)

var eventsCmd = &cobra.Command{
	Use:          "events",
	SilenceUsage: true,
	Short:        "Parse&&print upcoming events",
	RunE: func(cmd *cobra.Command, args []string) error {
		runtime, err := cfg.NewRuntime()
		if err != nil {
			return fmt.Errorf("create runtime: %w", err)
		}

		c, err := runtime.NewCalendar()
		if err != nil {
			return fmt.Errorf("create calendar: %w", err)
		}

		events, err := c.Events(context.Background(), calendar.DefaultLimit)
		if err != nil {
			return fmt.Errorf("fetch events: %w", err)
		}

		for _, e := range events {
			fmt.Printf(
				"[%s <--> %s] %s\n",
				e.Start.Format(time.RFC822), e.End.Format(time.RFC822),
				e.Summary,
			)
		}
		return nil
	},
}
