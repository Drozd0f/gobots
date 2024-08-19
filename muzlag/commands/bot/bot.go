package bot

import (
	"context"
	"fmt"

	"github.com/Drozd0f/gobots/muzlag/internal/bot"
	"github.com/Drozd0f/gobots/muzlag/internal/bot/handlers"
	"github.com/Drozd0f/gobots/muzlag/internal/config"
	"github.com/Drozd0f/gobots/muzlag/internal/service"
	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	pkgLog "github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

func RunBot(ctx context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config new config: %w", err)
	}

	logger := pkgLog.NewLogger(cfg.LogLevel)

	ph := handlers.NewPlayer(
		logger,
		service.NewService(
			logger,
			ytdl.NewDL(cfg.DL),
			ffmpeg.NewFfmpeg(cfg.Ffmpeg),
		))

	b, err := bot.NewBot(cfg, logger, ph)
	if err != nil {
		return fmt.Errorf("new bot: %w", err)
	}

	if err = b.Run(ctx); err != nil {
		return fmt.Errorf("run bot: %w", err)
	}

	return nil
}
