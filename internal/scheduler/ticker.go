package scheduler

import (
	"context"
	"time"
)

type Ticker struct {
	t *time.Ticker
}

type Event struct {
	Upcoming bool
	ToStart  time.Duration
	Left     time.Duration
	StartsAt time.Time
	EndsAt   time.Time
}

func (t *Ticker) Tick(ctx context.Context, d time.Duration) <-chan Event {
	ch := make(chan Event)
	go func() {
		for {
			toNextTick := time.Until(
				time.Now().Add(d).Truncate(d),
			)

			select {
			case <-ctx.Done():
				return
			case <-time.After(toNextTick):
			}
		}
	}()

	return ch
}
