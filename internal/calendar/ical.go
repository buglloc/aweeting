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

	ics "github.com/arran4/golang-ical"
	"github.com/buglloc/certifi"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const DefaultTimezone = "Local"

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

	parsed, err := ics.ParseCalendar(rsp.RawBody())
	if err != nil {
		return nil, fmt.Errorf("unable to parse calendar: %w", err)
	}

	minEnd := time.Now()
	maxEnd := minEnd.Add(limit)

	var events []Event
	for _, e := range parsed.Events() {
		eventID := e.Id()
		endAt, err := e.GetEndAt()
		if err != nil {
			log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid End datetime")
			continue
		}
		if endAt.Before(minEnd) {
			continue
		}

		startAt, err := e.GetStartAt()
		if err != nil {
			log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid Start datetime")
			continue
		}
		if startAt.After(maxEnd) {
			continue
		}

		if startAt.After(endAt) {
			log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid dates: Start after End")
			continue
		}

		var summary string
		if p := e.GetProperty(ics.ComponentPropertySummary); p != nil {
			summary = p.Value
		}

		events = append(events, Event{
			ID:      outEventID(summary, startAt.UTC().String(), endAt.UTC().String()),
			Summary: summary,
			Start:   startAt.In(c.loc),
			End:     endAt.In(c.loc),
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

func outEventID(summary, eventStart, eventEnd string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(eventStart))
	_, _ = h.Write([]byte(eventEnd))
	_, _ = h.Write([]byte(summary))
	return int(h.Sum32())
}
