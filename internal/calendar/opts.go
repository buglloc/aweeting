package calendar

import (
	"fmt"
	"time"
)

type Option func(c *ICal) error

func WithTimeZone(tz string) Option {
	return func(c *ICal) error {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return fmt.Errorf("invalid timezone: %w", err)
		}

		c.loc = loc
		return nil
	}
}
