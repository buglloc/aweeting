package calendar

import (
	"context"
	"crypto/tls"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/apognu/gocal"
	"github.com/buglloc/certifi"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type ICal struct {
	source string
	loc    *time.Location
	httpc  *resty.Client
}

func NewICal(source string, opts ...Option) (*ICal, error) {
	cal := &ICal{
		source: source,
		loc:    time.Local,
		httpc: resty.New().
			SetTLSClientConfig(&tls.Config{
				RootCAs: certifi.NewCertPool(),
			}).
			SetDoNotParseResponse(true).
			SetRetryCount(3).
			SetRetryWaitTime(100 * time.Millisecond).
			SetRetryMaxWaitTime(20 * time.Second).
			AddRetryCondition(func(rsp *resty.Response, err error) bool {
				return err != nil || rsp.StatusCode() == http.StatusInternalServerError
			}),
	}

	for _, opt := range opts {
		if err := opt(cal); err != nil {
			return nil, err
		}
	}

	return cal, nil
}

func (c *ICal) Events(ctx context.Context, limit time.Duration) ([]Event, error) {
	rsp, err := c.httpc.R().
		SetContext(ctx).
		Get(c.source)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch calendar: %w", err)
	}

	if rsp.IsError() {
		return nil, fmt.Errorf("non-200 response: %s", rsp.Status())
	}

	defer func() {
		_, _ = io.ReadAll(rsp.RawBody())
		_ = rsp.RawBody().Close()
	}()

	ical := gocal.NewParser(rsp.RawBody())
	start := time.Now()
	ical.Start = &start
	end := start.Add(limit)
	ical.End = &end
	ical.Strict = gocal.StrictParams{
		Mode: gocal.StrictModeFailAttribute,
	}

	if err := ical.Parse(); err != nil {
		return nil, fmt.Errorf("invalid ical events: %w", err)
	}

	events := make([]Event, 0, len(ical.Events))
	for _, e := range ical.Events {
		if e.Start == nil || e.End == nil {
			log.Warn().Any("event", e).Msg("skip event w/o boundaries")
			continue
		}

		if e.Summary == "" {
			e.Summary = "n/a"
		}

		events = append(events, Event{
			ID:      eventID(e.Summary, e.RawStart.Value, e.RawEnd.Value),
			Summary: e.Summary,
			Start:   e.Start.In(c.loc),
			End:     e.End.In(c.loc),
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	n := 0
	var lastEvent Event
	for _, e := range events {
		if lastEvent.IsSame(e) {
			continue
		}

		events[n] = e
		lastEvent = e
		n++
	}

	return events[:n], nil
}

func eventID(summary, eventStart, eventEnd string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(eventStart))
	_, _ = h.Write([]byte(eventEnd))
	_, _ = h.Write([]byte(summary))
	return int(h.Sum32())
}
