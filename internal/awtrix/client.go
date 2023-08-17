package awtrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/aweeting/internal/ticker"
)

const DefaultUpcomingLimit = 8 * time.Hour

var DefaultPayload = Payload{
	TextCase:    0,
	Color:       "#ffffff",
	Icon:        "11899",
	Repeat:      1,
	Duration:    5,
	Stack:       true,
	ScrollSpeed: 100,
}

type UpdaterConfig struct {
	Upstream        string
	Username        string
	Password        string
	Topic           string
	UpcomingLimit   time.Duration
	NonePayload     Payload
	UpcomingPayload Payload
	OnAirPayload    Payload
}

type MqttUpdater struct {
	mqtt mqtt.Client
	cfg  UpdaterConfig
}

func NewMqttUpdater(cfg UpdaterConfig) (*MqttUpdater, error) {
	if cfg.Topic == "" {
		return nil, errors.New(".Topic is required")
	}

	if cfg.Upstream == "" {
		return nil, errors.New(".Upstream is required")
	}

	l := log.With().Str("name", "awtrix.mqtt").Logger()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Upstream)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	opts.SetClientID("aweeting")
	opts.SetAutoReconnect(true)
	opts.OnConnect = func(_ mqtt.Client) {
		l.Info().Msg("connected")
	}
	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		l.Warn().Err(err).Msg("disconnected")
	}
	opts.OnReconnecting = func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		l.Info().Msg("reconnecting")
	}

	client := mqtt.NewClient(opts)
	client.Connect()

	return &MqttUpdater{
		mqtt: client,
		cfg:  cfg,
	}, nil
}

func (u *MqttUpdater) Update(ctx context.Context, event ticker.Event) error {
	payloadBytes, err := json.Marshal(u.payload(event))
	if err != nil {
		return fmt.Errorf("payload marshal: %w", err)
	}

	token := u.mqtt.Publish(u.cfg.Topic, 0, false, payloadBytes)
	select {
	case <-token.Done():
		return token.Error()
	case <-ctx.Done():
		return fmt.Errorf("canceled: %w", ctx.Err())
	}
}

func (u *MqttUpdater) payload(event ticker.Event) Payload {
	var payload Payload
	switch {
	case event.IsZero() || event.StartsAt.Sub(time.Now()) > u.cfg.UpcomingLimit:
		payload = u.cfg.NonePayload
	case event.Upcoming:
		payload = u.cfg.UpcomingPayload
	default:
		payload = u.cfg.OnAirPayload
	}

	payload.Text = u.eventText(event)
	return payload
}

func (u *MqttUpdater) eventText(event ticker.Event) string {
	switch {
	case event.IsZero():
		return " ##:##"
	case event.Upcoming:
		return fmt.Sprintf("-%s", u.formatDuration(event.ToStart))
	default:
		return fmt.Sprintf(" %s", u.formatDuration(event.Left))
	}
}

func (u *MqttUpdater) formatDuration(d time.Duration) string {
	if d.Minutes() < 60.0 {
		return fmt.Sprintf("00:%02d", int(d.Minutes()))
	}

	if d.Hours() < 24.0 {
		remainingMinutes := math.Mod(d.Minutes(), 60)
		return fmt.Sprintf("%02d:%02d", int(d.Hours()), int(remainingMinutes))
	}

	return "##:##"
}
