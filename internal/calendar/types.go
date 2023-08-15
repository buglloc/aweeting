package calendar

import "time"

type Event struct {
	ID      int
	Summary string
	Start   time.Time
	End     time.Time
}

func (e *Event) IsSame(other Event) bool {
	if e.ID != "" && other.ID != "" {
		return e.ID == other.ID
	}

	return e.Summary == other.Summary &&
		e.Start.Compare(other.Start) == 0 &&
		e.End.Compare(other.End) == 0
}
