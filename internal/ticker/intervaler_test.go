package ticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/buglloc/aweeting/internal/calendar"
)

var now time.Time

func init() {
	now = time.Unix(544672800, 0)
	nowFn = func() time.Time {
		return now
	}
}

func TestIntervaler_expired(t *testing.T) {
	cases := []struct {
		name     string
		events   []calendar.Event
		expected Interval
	}{
		{
			name: "one",
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(-100 * time.Minute),
					End:   now.Add(-50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(20 * time.Minute),
					End:   now.Add(30 * time.Minute),
				},
				{
					ID:    3,
					Start: now.Add(40 * time.Minute),
					End:   now.Add(50 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(20 * time.Minute),
				End:   now.Add(30 * time.Minute),
			},
		},
		{
			name: "multiple",
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(-100 * time.Minute),
					End:   now.Add(-50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(-20 * time.Minute),
					End:   now.Add(-30 * time.Minute),
				},
				{
					ID:    3,
					Start: now.Add(40 * time.Minute),
					End:   now.Add(50 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(40 * time.Minute),
				End:   now.Add(50 * time.Minute),
			},
		},
		{
			name: "zero",
			events: []calendar.Event{
				{
					ID:    3,
					Start: now.Add(-400 * time.Minute),
					End:   now.Add(-500 * time.Minute),
				},
				{
					ID:    1,
					Start: now.Add(-100 * time.Minute),
					End:   now.Add(-50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(-30 * time.Minute),
					End:   now.Add(-20 * time.Minute),
				},
			},
			expected: Interval{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			i := NewIntervaler(0)
			i.UpdateEvents(tc.events)

			actual := i.Current()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestIntervaler_overlaps(t *testing.T) {
	cases := []struct {
		name     string
		events   []calendar.Event
		expected Interval
	}{
		{
			name: "full-overlap",
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(100 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(20 * time.Minute),
					End:   now.Add(30 * time.Minute),
				},
				{
					ID:    3,
					Start: now.Add(50 * time.Minute),
					End:   now.Add(100 * time.Minute),
				},
				{
					ID:    4,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(100 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(10 * time.Minute),
				End:   now.Add(100 * time.Minute),
			},
		},
		{
			name: "partial-overlap",
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(100 * time.Minute),
				},
				{
					ID:    3,
					Start: now.Add(100 * time.Minute),
					End:   now.Add(110 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(10 * time.Minute),
				End:   now.Add(100 * time.Minute),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			i := NewIntervaler(0)
			i.UpdateEvents(tc.events)

			actual := i.Current()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestIntervaler_merge(t *testing.T) {
	cases := []struct {
		name     string
		jitter   time.Duration
		events   []calendar.Event
		expected Interval
	}{
		{
			name:   "split-no-jitter",
			jitter: 0,
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(50 * time.Minute),
					End:   now.Add(60 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(10 * time.Minute),
				End:   now.Add(50 * time.Minute),
			},
		},
		{
			name:   "merge-no-jitter",
			jitter: 0,
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(49 * time.Minute),
					End:   now.Add(60 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(10 * time.Minute),
				End:   now.Add(60 * time.Minute),
			},
		},
		{
			name:   "merge-w-jitter",
			jitter: 1 * time.Minute,
			events: []calendar.Event{
				{
					ID:    1,
					Start: now.Add(10 * time.Minute),
					End:   now.Add(50 * time.Minute),
				},
				{
					ID:    2,
					Start: now.Add(50 * time.Minute),
					End:   now.Add(60 * time.Minute),
				},
			},
			expected: Interval{
				Start: now.Add(10 * time.Minute),
				End:   now.Add(60 * time.Minute),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			i := NewIntervaler(tc.jitter)
			i.UpdateEvents(tc.events)

			actual := i.Current()
			require.Equal(t, tc.expected, actual)
		})
	}
}
