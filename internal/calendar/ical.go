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
	"github.com/teambition/rrule-go"
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

	now := time.Now()
	tb := TimeBound{
		Start: now,
		End:   now.Add(limit),
	}
	var events []Event
	for _, e := range parsed.Events() {
		var summary string
		if p := e.GetProperty(ics.ComponentPropertySummary); p != nil {
			summary = p.Value
		}

		for _, times := range c.eventTimes(e, tb) {
			events = append(events, Event{
				ID:      outEventID(summary, times.Start.UTC().String(), times.End.UTC().String()),
				Summary: summary,
				Start:   times.Start.In(c.loc),
				End:     times.End.In(c.loc),
			})
		}
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

func (c *ICal) eventTimes(e *ics.VEvent, tb TimeBound) []TimeBound {
	if prop := e.GetProperty(ics.ComponentPropertyRrule); prop != nil {
		return c.rrEventTimes(e, prop, tb)
	}

	eventID := e.Id()
	endAt, err := e.GetEndAt()
	if err != nil {
		log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid End datetime")
		return nil
	}

	if endAt.Before(tb.Start) {
		return nil
	}

	startAt, err := e.GetStartAt()
	if err != nil {
		log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid Start datetime")
		return nil
	}

	if startAt.After(tb.End) {
		return nil
	}

	if startAt.After(endAt) {
		log.Warn().Str("event_id", eventID).Err(err).Msg("skip event with invalid dates: Start after End")
		return nil
	}

	return []TimeBound{{Start: startAt, End: endAt}}
}

func (c *ICal) rrEventTimes(e *ics.VEvent, rrProp *ics.IANAProperty, tb TimeBound) []TimeBound {
	eventID := e.Id()

	duration, err := c.eventDuration(e)
	if err != nil {
		log.Warn().Str("event_id", eventID).Err(err).Msg("skip recurring event with invalid duration")
		return nil
	}

	startProp := e.GetProperty(ics.ComponentPropertyDtStart)
	if startProp == nil {
		log.Warn().Str("event_id", eventID).Msg("skip recurring event w/o dstart")
		return nil
	}

	var tzStr string
	if tzID := startProp.ICalParameters["TZID"]; len(tzID) > 0 {
		tzStr = fmt.Sprintf("TZID=%s:", tzID[0])
	}

	rfcRRStr := fmt.Sprintf("DTSTART:%s%s\nRRULE:%s", tzStr, startProp.Value, rrProp.Value)
	rOption, err := rrule.StrToROptionInLocation(rfcRRStr, c.loc)
	if err != nil {
		log.Warn().
			Str("event_id", eventID).
			Str("rrule", rfcRRStr).
			Err(err).
			Msg("skip recurring event with invalid rrule")
		return nil
	}

	rr, err := rrule.NewRRule(*rOption)
	if err != nil {
		log.Warn().
			Str("event_id", eventID).
			Str("rrule", rfcRRStr).
			Err(err).
			Msg("skip recurring event with invalid rrule")
		return nil
	}

	var out []TimeBound
	for _, start := range rr.Between(tb.Start, tb.End, true) {
		end := start.Add(duration)
		out = append(out, TimeBound{
			Start: start,
			End:   end,
		})
	}

	return out
}

func (c *ICal) eventDuration(e *ics.VEvent) (time.Duration, error) {
	startAt, err := e.GetStartAt()
	if err != nil {
		return 0, fmt.Errorf("invalid start at: %w", err)
	}

	endAt, err := e.GetEndAt()
	if err != nil {
		return 0, fmt.Errorf("invalid end at: %w", err)
	}

	d := endAt.Sub(startAt)
	if d <= 0 {
		return 0, fmt.Errorf("invalid duration: %s -> %s", startAt, endAt)
	}

	return d, nil
}
