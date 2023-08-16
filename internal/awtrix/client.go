package awtrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buglloc/aweeting/internal/ticker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"math"
	"time"
)

type ColorSet struct {
	Zero     []int
	Upcoming []int
	OnAir    []int
}

type IconSet struct {
	Zero     string
	Upcoming string
	OnAir    string
}

type MqttConfig struct {
	Upstream string
	Username string
	Password string
	Topic    string
	Icons    IconSet
	Colors   ColorSet
}

type MqttUpdater struct {
	mqtt   mqtt.Client
	topic  string
	icons  IconSet
	colors ColorSet
}

func NewMqttUpdater(cfg MqttConfig) (*MqttUpdater, error) {
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
		mqtt:   client,
		topic:  cfg.Topic,
		icons:  cfg.Icons,
		colors: cfg.Colors,
	}, nil
}

func (u *MqttUpdater) Update(ctx context.Context, event ticker.Event) error {
	payload := Payload{
		Text:  u.eventText(event),
		Color: u.eventTextColor(event),
		Icon:  u.eventIcon(event),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload marshal: %w", err)
	}

	token := u.mqtt.Publish(u.topic, 0, false, payloadBytes)
	select {
	case <-token.Done():
		return token.Error()
	case <-ctx.Done():
		return fmt.Errorf("canceled: %w", ctx.Err())
	}
}

func (u *MqttUpdater) eventText(event ticker.Event) string {
	switch {
	case event.IsZero():
		return " ##:##"
	case event.Upcoming:
		return fmt.Sprintf("-%s", formatDuration(event.ToStart))
	default:
		return fmt.Sprintf(" %s", formatDuration(event.Left))
	}
}

func (u *MqttUpdater) eventTextColor(event ticker.Event) []int {
	switch {
	case event.IsZero():
		return u.colors.Zero
	case event.Upcoming:
		return u.colors.Upcoming
	default:
		return u.colors.OnAir
	}
}

func (u *MqttUpdater) eventIcon(event ticker.Event) string {
	switch {
	case event.IsZero():
		return u.icons.Zero
	case event.Upcoming:
		return u.icons.Upcoming
	default:
		return u.icons.OnAir
	}
}

func formatDuration(d time.Duration) string {
	if d.Minutes() < 60.0 {
		return fmt.Sprintf("00:%02d", int(d.Minutes()))
	}

	if d.Hours() < 24.0 {
		remainingMinutes := math.Mod(d.Minutes(), 60)
		return fmt.Sprintf("%02d:%02d", int(d.Hours()), int(remainingMinutes))
	}

	return "##:##"
}
