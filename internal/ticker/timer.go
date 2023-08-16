package ticker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type TimerConfig struct {
	Name     string
	Interval time.Duration
}

type Timer struct {
	fn       func(ctx context.Context) error
	name     string
	interval time.Duration
	timer    *time.Timer
}

func NewTimer(fn func(ctx context.Context) error, cfg TimerConfig) *Timer {
	interval := 1 * time.Minute
	if cfg.Interval > 0 {
		interval = cfg.Interval
	}

	return &Timer{
		fn:       fn,
		name:     cfg.Name,
		interval: interval,
	}
}

func (t *Timer) Start(ctx context.Context) {
	tick := func() time.Duration {
		return time.Until(
			time.Now().Add(t.interval).Truncate(t.interval),
		)
	}

	ctx = log.With().Str("name", t.name).Logger().WithContext(context.Background())
	t.timer = time.AfterFunc(tick(), func() {
		now := time.Now()
		log.Ctx(ctx).Info().Msg("task started")
		defer log.Ctx(ctx).Info().Dur("elapsed", time.Since(now)).Msg("task finished")

		if ctx.Err() != nil {
			log.Ctx(ctx).Info().Msg("task stopped")
			return
		}

		if err := t.fn(ctx); err != nil {
			log.Ctx(ctx).Err(err).Msg("task failed")
		}

		t.timer.Reset(tick())
	})

	go func() {
		<-ctx.Done()
		t.timer.Stop()
	}()
}
