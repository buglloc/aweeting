package ticker

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/buglloc/aweeting/internal/calendar"
)

const (
	DefaultJitter        = 20 * time.Minute
	DefaultPreviewLimit  = 24 * time.Hour
	DefaultFetchInterval = 1 * time.Hour
	DefaultTickInterval  = 5 * time.Minute
)

var _ Ticker = (*ConstTicker)(nil)

type ConstTickerConfig struct {
	Jitter        time.Duration
	PreviewLimit  time.Duration
	FetchInterval time.Duration
	TickInterval  time.Duration
}

type ConstTicker struct {
	cal           calendar.Calendar
	ctx           context.Context
	cancelCtx     context.CancelFunc
	done          chan struct{}
	interval      *Intervaler
	previewLimit  time.Duration
	fetchInterval time.Duration
	tickInterval  time.Duration
}

func NewConstTicker(cal calendar.Calendar, cfg ConstTickerConfig) (*ConstTicker, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConstTicker{
		cal:           cal,
		ctx:           ctx,
		cancelCtx:     cancel,
		done:          make(chan struct{}),
		interval:      NewIntervaler(cfg.Jitter),
		previewLimit:  cfg.PreviewLimit,
		fetchInterval: cfg.FetchInterval,
		tickInterval:  cfg.TickInterval,
	}, nil
}

func (t *ConstTicker) Start(handler Handler) error {
	defer close(t.done)

	if err := t.fetchEvents(t.ctx); err != nil {
		return fmt.Errorf("first update events: %w", err)
	}

	handle := t.newTickHandle(handler)
	if err := handle(t.ctx); err != nil {
		return fmt.Errorf("first tick: %w", err)
	}

	log.Info().Msg("const ticker started")
	NewTimer(
		t.fetchEvents,
		TimerConfig{
			Name:     "fetch",
			Interval: t.fetchInterval,
		},
	).Start(t.ctx)

	NewTimer(
		handle,
		TimerConfig{
			Name:     "tick",
			Interval: t.tickInterval,
		},
	).Start(t.ctx)

	<-t.ctx.Done()
	return nil
}

func (t *ConstTicker) Stop(ctx context.Context) {
	t.cancelCtx()
	select {
	case <-ctx.Done():
		return
	case <-t.done:
		return
	}
}

func (t *ConstTicker) fetchEvents(ctx context.Context) error {
	events, err := t.cal.Events(ctx, t.previewLimit)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}
	log.Ctx(ctx).Info().Int("count", len(events)).Msg("got calendar events")

	t.interval.UpdateEvents(events)
	return nil
}

func (t *ConstTicker) newTickHandle(handler Handler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		event := t.interval.Current().ToEvent(time.Now().Truncate(time.Minute))
		return handler(ctx, event)
	}
}
