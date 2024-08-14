package config

import (
	"fmt"
	"log/slog"

	"github.com/Drozd0f/gobots/muzlag/pkg/config"
	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

const AppName = "muzlag"

type Config struct {
	Token    string        `required:"true"`
	Prefix   string        `required:"true"`
	LogLevel slog.Level    `required:"true" split_words:"true"`
	DL       ytdl.Config   `required:"true"`
	Ffmpeg   ffmpeg.Config `required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg = new(Config)
	if err := config.Process(AppName, cfg); err != nil {
		return nil, fmt.Errorf("process config: %w", err)
	}

	return cfg, nil
}
