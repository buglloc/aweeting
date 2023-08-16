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
	activeEvent  calendar.Event
	gap          time.Duration
	updInterval  time.Duration
	previewLimit time.Duration
}

func NewScheduler(cal calendar.Calendar) *Scheduler {
	sched := quartz.NewStdSchedulerWithOptions(quartz.StdSchedulerOptions{
		BlockingExecution: true,
	})
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

	updateTick := func() time.Duration {
		return time.Until(
			time.Now().Add(s.updInterval).Truncate(s.updInterval),
		)
	}

	s.sched.Start(s.ctx)
	toUpdateTick := updateTick()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(toUpdateTick):
			toUpdateTick = updateTick()
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

	activeIDs := make(map[int]struct{})
	for _, id := range s.sched.GetJobKeys() {
		activeIDs[id] = struct{}{}
	}

	for _, e := range events {
		if _, ok := activeIDs[e.ID]; ok {
			continue
		}

		s.sched.ScheduleJob(s.ctx)
	}
	s.sched.GetJobKeys()
	time.AfterFunc()
	ss := time.NewTimer()
	ss.Reset()
}
