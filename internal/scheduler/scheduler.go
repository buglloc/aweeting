package scheduler

import (
	"context"
	"fmt"
	"github.com/buglloc/aweeting/internal/calendar"
	"github.com/reugn/go-quartz/quartz"
	"github.com/rs/zerolog/log"
	"time"
)

type Scheduler struct {
	cal          calendar.Calendar
	sched        quartz.Scheduler
	ctx          context.Context
	cancelCtx    context.CancelFunc
	done         chan struct{}
	gap          time.Duration
	jitter       time.Duration
	updInterval  time.Duration
	previewLimit time.Duration
}

func NewScheduler(cal calendar.Calendar) *Scheduler {
	sched := quartz.NewStdSchedulerWithOptions(quartz.StdSchedulerOptions{})
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cal:          cal,
		sched:        sched,
		ctx:          ctx,
		cancelCtx:    cancel,
		done:         make(chan struct{}),
		gap:          10 * time.Minute,
		jitter:       20 * time.Minute,
		updInterval:  1 * time.Hour,
		previewLimit: 24 * time.Hour,
	}
}

func (s *Scheduler) Start() error {
	defer close(s.done)

	s.sched.Start(s.ctx)
	for {
		toNextTick := time.Until(
			time.Now().Add(s.updInterval).Truncate(s.updInterval),
		)

		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(toNextTick):
			if err := s.updateJobs(); err != nil {

			}
		}
	}
}

func (s *Scheduler) Stop(ctx context.Context) {
	s.cancelCtx()
	select {
	case <-ctx.Done():
		return
	case <-s.done:
		return
	}
}

func (s *Scheduler) updateJobs() error {
	events, err := s.cal.Events(s.ctx, s.previewLimit)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}
	log.Info().Int("count", len(events)).Msg("got calendar events")

}
