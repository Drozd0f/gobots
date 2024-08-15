package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/bot/handlers"
	"github.com/Drozd0f/gobots/muzlag/internal/config"
)

type Bot struct {
	logger  *slog.Logger
	session *discordgo.Session
}

func NewBot(cfg config.Config, logger *slog.Logger, ph handlers.Player) (Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return Bot{}, fmt.Errorf("error creating discord session: %w", err)
	}

	// Register intents
	session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMembers

	// Setup state
	session.StateEnabled = true
	session.State.TrackChannels = true
	session.State.TrackMembers = true
	session.State.TrackVoice = true
	session.State.TrackPresences = true

	registerRoutes(session, newRouterParams{
		Logger: logger,
		Prefix: cfg.Prefix,
		Player: ph,
	})

	return Bot{
		logger:  logger,
		session: session,
	}, nil
}

func (b Bot) Run(ctx context.Context) error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("error opening session: %w", err)
	}

	b.logger.Info("Bot is now running.  Press CTRL-C to exit.", slog.String("name", config.AppName))

	<-ctx.Done()

	if err := b.session.Close(); err != nil {
		return fmt.Errorf("error closing session: %w", err)
	}

	b.logger.Info("Bot stopped successfully.", slog.String("name", config.AppName))

	return nil
}
