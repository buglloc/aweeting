package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/buglloc/aweeting/internal/awtrix"
	"github.com/buglloc/aweeting/internal/calendar"
	"github.com/buglloc/aweeting/internal/ticker"
)

type Config struct {
	Verbose  bool     `koanf:"verbose"`
	Calendar Calendar `koanf:"calendar"`
	Ticker   Ticker   `koanf:"ticker"`
	Mqtt     Mqtt     `koanf:"mqtt"`
	Awtrix   Awtrix   `koanf:"awtrix"`
}

func (c *Config) Validate() error {
	return nil
}

type Runtime struct {
	cfg *Config
}

func LoadConfig(files ...string) (*Config, error) {
	out := Config{
		Calendar: Calendar{
			Timezone: calendar.DefaultTimezone,
		},
		Ticker: Ticker{
			Jitter:        ticker.DefaultJitter,
			PreviewLimit:  ticker.DefaultPreviewLimit,
			FetchInterval: ticker.DefaultFetchInterval,
			TickInterval:  ticker.DefaultTickInterval,
		},
		Awtrix: Awtrix{
			UpcomingLimit: awtrix.DefaultUpcomingLimit,
			Messages: AwtrixMessagesSet{
				None:     AwtrixMessage(awtrix.DefaultPayload),
				Upcoming: AwtrixMessage(awtrix.DefaultPayload),
				OnAir:    AwtrixMessage(awtrix.DefaultPayload),
			},
		},
	}

	k := koanf.New(".")

	yamlParser := yaml.Parser()
	for _, fpath := range files {
		if err := k.Load(file.Provider(fpath), yamlParser); err != nil {
			return nil, fmt.Errorf("load %q config: %w", fpath, err)
		}
	}

	envCb := func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "AW_")), "_", ".", -1)
	}
	if err := k.Load(env.Provider("AW_", ".", envCb), nil); err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	return &out, k.Unmarshal("", &out)
}

func (c *Config) NewRuntime() (*Runtime, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &Runtime{
		cfg: c,
	}, nil
}
