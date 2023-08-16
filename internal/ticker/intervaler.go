package ticker

import (
	"fmt"
	"github.com/buglloc/aweeting/internal/calendar"
	"sync"
	"time"
)

type Intervaler struct {
	mu     sync.RWMutex
	events []calendar.Event
	jitter time.Duration
}

type Interval struct {
	Start time.Time
	End   time.Time
}

func NewIntervaler(jitter time.Duration) *Intervaler {
	return &Intervaler{
		jitter: jitter,
	}
}

func (c *Intervaler) UpdateEvents(events []calendar.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.events = events
}

func (c *Intervaler) Current() Interval {
	c.mu.RLock()
	defer c.mu.RUnlock()

	skip := 0
	now := nowFn()
	for _, e := range c.events {
		if now.Before(e.End) {
			break
		}

		skip++
	}
	c.events = c.events[skip:]

	if len(c.events) == 0 {
		return Interval{}
	}

	cur := Interval{
		Start: c.events[0].Start,
		End:   c.events[0].End,
	}

	for _, e := range c.events[1:] {
		switch {
		case cur.Start.Before(e.Start) && cur.End.After(e.End):
			// overlap
			continue
		case cur.End.Add(c.jitter).After(e.Start):
			cur.End = e.End
			continue
		}

		break
	}

	return cur
}

func (i Interval) String() string {
	return fmt.Sprintf("%s -> %s", i.Start.Format(time.RFC3339), i.End.Format(time.RFC3339))
}

func (i Interval) IsZero() bool {
	return i.Start.IsZero()
}

func (i Interval) ToEvent(now time.Time) Event {
	if i.Start.IsZero() {
		return Event{
			Upcoming: true,
		}
	}

	return Event{
		Upcoming: i.Start.After(now),
		ToStart:  i.Start.Sub(now),
		Left:     i.End.Sub(now),
		StartsAt: i.Start,
		EndsAt:   i.End,
	}
}
