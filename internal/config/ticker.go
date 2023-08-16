package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/buglloc/aweeting/internal/ticker"
)

type Ticker struct {
	Jitter        time.Duration `koanf:"jitter"`
	PreviewLimit  time.Duration `koanf:"previewLimit"`
	FetchInterval time.Duration `koanf:"fetchInterval"`
	TickInterval  time.Duration `koanf:"tickInterval"`
}

func (c *Ticker) Validate() error {
	if c.PreviewLimit == 0 {
		return errors.New(".PreviewLimit is required")
	}

	if c.FetchInterval == 0 {
		return errors.New(".FetchInterval is required")
	}

	if c.TickInterval == 0 {
		return errors.New(".TickInterval is required")
	}

	return nil
}

func (r *Runtime) NewTicker() (ticker.Ticker, error) {
	if err := r.cfg.Calendar.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cal, err := r.NewCalendar()
	if err != nil {
		return nil, fmt.Errorf("create calendar: %w", err)
	}

	return ticker.NewConstTicker(cal, ticker.ConstTickerConfig(r.cfg.Ticker))
}
