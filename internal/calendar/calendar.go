package calendar

import (
	"context"
	"time"
)

const DefaultLimit = 3 * 24 * time.Hour

type Calendar interface {
	Events(ctx context.Context, limit time.Duration) ([]Event, error)
}
