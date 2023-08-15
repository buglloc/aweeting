package notifier

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"a.yandex-team.ru/junk/buglloc/alice-notifier/internal/bass"
	"a.yandex-team.ru/junk/buglloc/alice-notifier/internal/calendar"
	"a.yandex-team.ru/junk/buglloc/alice-notifier/internal/config"
	"a.yandex-team.ru/junk/buglloc/alice-notifier/internal/logger"
	"a.yandex-team.ru/library/go/core/log"
	"github.com/karlseguin/ccache/v2"
)

const cacheTTL = 30 * time.Minute

type Notifier struct {
	cfg         *config.Config
	cc          *calendar.Calendar
	bassc       *bass.Bass
	excludes    map[string]struct{}
	eventsCache *ccache.Cache
	ctx         context.Context
	cancelCtx   context.CancelFunc
	closed      chan struct{}
}

func NewNotifier(cfg *config.Config) (*Notifier, error) {
	excludes := make(map[string]struct{}, len(cfg.ExcludedEvents))
	for _, e := range cfg.ExcludedEvents {
		excludes[e] = struct{}{}
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	return &Notifier{
		cfg:      cfg,
		cc:       calendar.NewCalendar(cfg.Calendar.Upstream, cfg.Calendar.Token),
		bassc:    bass.NewBass(cfg.Bass.Endpoint),
		excludes: excludes,
		eventsCache: ccache.New(
			ccache.Configure().MaxSize(1024),
		),
		ctx:       ctx,
		cancelCtx: cancelCtx,
		closed:    make(chan struct{}),
	}, nil
}

func (n *Notifier) Start() error {
	defer close(n.closed)

	for {
		toNextWork := time.Until(
			time.Now().Add(n.cfg.RunEvery).Truncate(n.cfg.RunEvery),
		)
		t := time.NewTimer(toNextWork)

		select {
		case <-n.ctx.Done():
			t.Stop()
			return nil
		case <-t.C:
			t.Stop()

			logger.Info("process events")
			if err := n.process(n.ctx); err != nil {
				logger.Error("oops, unable to process", log.Error(err))
			} else {
				logger.Info("events processed")
			}
		}
	}
}

func (n *Notifier) Shutdown(ctx context.Context) {
	n.cancelCtx()

	select {
	case <-ctx.Done():
	case <-n.closed:
	}

	n.eventsCache.Stop()
}

func (n *Notifier) process(ctx context.Context) error {
	events, err := n.cc.Events(ctx)
	if err != nil {
		return fmt.Errorf("unable to get calendar evetns: %w", err)
	}

	now := time.Now()
	var nearestTime time.Duration
	var subjects []string
	for _, e := range events {
		if _, ok := n.excludes[e.Subject]; ok {
			continue
		}

		cacheKey := fmt.Sprintf("%s_%s", e.StartTime, e.Subject)
		if ce := n.eventsCache.Get(cacheKey); ce != nil {
			continue
		}

		left := e.StartTime.Sub(now)
		if left > n.cfg.Gap {
			continue
		}

		if nearestTime == 0 || left < nearestTime {
			nearestTime = left
		}

		subjects = append(subjects, e.Subject)
		n.eventsCache.Set(cacheKey, struct{}{}, cacheTTL)
	}

	if len(subjects) == 0 {
		return nil
	}

	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("Андрей, у тебя ")
	if len(subjects) > 1 {
		msgBuilder.WriteString(strconv.Itoa(len(subjects)))
		msgBuilder.WriteString(" встречки через ")
	} else {
		msgBuilder.WriteString(" встречка через ")
	}

	msgBuilder.WriteString(strconv.Itoa(int(nearestTime.Minutes())))
	msgBuilder.WriteString(" минут ")

	for i, s := range subjects {
		if i != 0 {
			msgBuilder.WriteString(" и ")
		}
		msgBuilder.WriteByte('"')
		msgBuilder.WriteString(s)
		msgBuilder.WriteByte('"')
	}

	msg := msgBuilder.String()
	logger.Info("notify",
		log.String("msg", msg),
		log.String("device_id", n.cfg.DeviceID),
		log.String("user_id", n.cfg.UserID),
	)

	return n.bassc.Say(
		ctx,
		n.cfg.DeviceID,
		n.cfg.UserID,
		msg,
	)
}
