package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/buglloc/aweeting/internal/awtrix"
)

type Mqtt struct {
	Upstream string `koanf:"upstream"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Topic    string `koanf:"topic"`
}

type Awtrix struct {
	UpcomingLimit time.Duration     `koanf:"upcomingLimit"`
	Messages      AwtrixMessagesSet `koanf:"messages"`
}

type AwtrixMessagesSet struct {
	None     AwtrixMessage `koanf:"none"`
	Upcoming AwtrixMessage `koanf:"upcoming"`
	OnAir    AwtrixMessage `koanf:"onAir"`
}

type AwtrixMessage struct {
	// The text to display
	Text string `koanf:"text"`
	// Changes the Uppercase setting. 0=global setting, 1=forces uppercase; 2=shows as it sent
	TextCase int `koanf:"textCase"`
	// Draw the text on top
	TopText bool `koanf:"topText"`
	//Sets an offset for the x position of a starting text
	TextOffset int `koanf:"textOffset"`
	// The text, bar or line color (#hex)
	Color string `koanf:"color"`
	// Sets a background color (#hex)
	Background string `koanf:"background"`
	// Fades each letter in the text differently through the entire RGB spectrum
	Rainbow bool `koanf:"rainbow"`
	// The icon ID or filename (without extension) to display on the app
	Icon string `koanf:"icon"`
	// 0 = Icon doesn't move. 1 = Icon moves with text and will not appear again. 2 = Icon moves with text but appears again when the text starts to scroll again.
	PushIcon int `koanf:"pushIcon"`
	// Sets how many times the text should be scrolled through the matrix before the app ends
	Repeat int `koanf:"repeat"`
	// Sets how long the app or notification should be displayed
	Duration int `koanf:"duration"`
	// Enables or disables autoscaling for bar and linechart
	Autoscale bool `koanf:"autoscale"`
	//  Defines the position of your custompage in the loop, starting at 0 for the first position. This will only apply with your first push. This function is experimental
	Pos int `koanf:"pos"`
	// Removes the custom app when there is no update after the given time in seconds
	Lifetime int `koanf:"lifetime"`
	// Defines if the **notification** will be stacked. false will immediately replace the current notification
	Stack bool `koanf:"stack"`
	// If the Matrix is off, the notification will wake it up for the time of the notification
	Wakeup bool `koanf:"wakeup"`
	// Disables the textscrolling
	NoScroll bool `koanf:"noScroll"`
	// Modifies the scrollspeed. You need to enter a percentage value
	ScrollSpeed int `koanf:"scrollSpeed"`
	// Shows an (https://blueforcer.github.io/awtrix-light/#/effects) as background
	Effect string `koanf:"effect"`
	// Changes color and speed of the (https://blueforcer.github.io/awtrix-light/#/effects)
	EffectSettings map[string]any `koanf:"effectSettings"`
}

func (c *Mqtt) Validate() error {
	if c.Upstream == "" {
		return errors.New(".Upstream is required")
	}

	if c.Topic == "" {
		return errors.New(".Topic is required")
	}

	return nil
}

func (r *Runtime) NewAwtrixUpdater() (*awtrix.MqttUpdater, error) {
	if err := r.cfg.Mqtt.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mqtt config: %w", err)
	}

	return awtrix.NewMqttUpdater(awtrix.UpdaterConfig{
		Upstream:        r.cfg.Mqtt.Upstream,
		Username:        r.cfg.Mqtt.Username,
		Password:        r.cfg.Mqtt.Password,
		Topic:           r.cfg.Mqtt.Topic,
		UpcomingLimit:   r.cfg.Awtrix.UpcomingLimit,
		NonePayload:     awtrix.Payload(r.cfg.Awtrix.Messages.None),
		UpcomingPayload: awtrix.Payload(r.cfg.Awtrix.Messages.Upcoming),
		OnAirPayload:    awtrix.Payload(r.cfg.Awtrix.Messages.OnAir),
	})
}
