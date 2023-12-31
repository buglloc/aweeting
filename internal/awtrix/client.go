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

const (
	DefaultUpcomingLimit  = 8 * time.Hour
	MqttConnectionTimeout = 5 * time.Minute
)

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
	SelfDestruct    bool
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
	if token := client.Connect(); token.WaitTimeout(MqttConnectionTimeout) && token.Error() != nil {
		return nil, fmt.Errorf("MQTT connection failed: %w", token.Error())
	}

	return &MqttUpdater{
		mqtt: client,
		cfg:  cfg,
	}, nil
}

func (u *MqttUpdater) Update(ctx context.Context, event ticker.Event) error {
	payloadBytes, err := u.payloadBytes(event)
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

func (u *MqttUpdater) payloadBytes(event ticker.Event) ([]byte, error) {
	var payload Payload
	switch {
	case u.isNoneEvent(event):
		if u.cfg.SelfDestruct {
			return nil, nil
		}

		payload = u.cfg.NonePayload
	case event.Upcoming:
		payload = u.cfg.UpcomingPayload
	default:
		payload = u.cfg.OnAirPayload
	}

	payload.Text = u.eventText(event)
	return json.Marshal(payload)
}

func (u *MqttUpdater) eventText(event ticker.Event) string {
	switch {
	case u.isNoneEvent(event):
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
func (u *MqttUpdater) isNoneEvent(event ticker.Event) bool {
	return event.IsZero() || time.Until(event.StartsAt) > u.cfg.UpcomingLimit
}
