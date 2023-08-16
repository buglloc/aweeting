package ticker

import (
	"context"
	"time"
)

type Ticker interface {
	Start(Handler) error
	Stop(context.Context)
}

type Handler func(ctx context.Context, event Event) error

type Event struct {
	Upcoming bool
	ToStart  time.Duration
	Left     time.Duration
	StartsAt time.Time
	EndsAt   time.Time
}

func (e *Event) IsZero() bool {
	return e.StartsAt.IsZero()
}
