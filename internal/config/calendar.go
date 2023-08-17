package config

import (
	"errors"
	"fmt"

	"github.com/buglloc/aweeting/internal/calendar"
)

type Calendar struct {
	SourceURL string `koanf:"sourceUrl"`
	Timezone  string `koanf:"timezone"`
}

func (c *Calendar) Validate() error {
	if c.SourceURL == "" {
		return errors.New(".SourceURL is required")
	}

	return nil
}

func (r *Runtime) NewCalendar() (calendar.Calendar, error) {
	if err := r.cfg.Calendar.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return calendar.NewICal(r.cfg.Calendar.SourceURL,
		calendar.WithTimeZone(r.cfg.Calendar.Timezone),
	)
}
