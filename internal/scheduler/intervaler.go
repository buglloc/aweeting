package scheduler

import (
	"github.com/buglloc/aweeting/internal/calendar"
	"time"
)

type IntervalCalculator struct {
	gap    time.Duration
	jitter time.Duration
}

type Interval struct {
	Start, End time.Time
}

func NewIntervalCalculator(gap, jitter time.Duration) *IntervalCalculator {
	return &IntervalCalculator{
		gap:    gap,
		jitter: jitter,
	}
}

func (c *IntervalCalculator) Calculate(events ...calendar.Event) []Interval {
	if len(events) == 0 {
		return nil
	}

	intervals := make([]Interval, 0, len(events))
	ci := Interval{
		Start: events[0].Start,
		End:   events[0].End,
	}
	for _, ni := range events[1:] {

	}
}

func (i *Interval) Overlaps(o Interval) bool {
	return i.Start.Compare(o.Start) <= 0 && i.End.Compare(o.End) <= 0
}
