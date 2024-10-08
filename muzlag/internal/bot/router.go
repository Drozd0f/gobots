package bot

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/bot/handlers"
	"github.com/Drozd0f/gobots/muzlag/internal/bot/middleware"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
)

type newRouterParams struct {
	Prefix string
	Player handlers.Player
	Logger *slog.Logger
}

func registerRoutes(session *discordgo.Session, p newRouterParams) {
	registerMessageCreateHandlers(session, registerMessageCreateHandlersParams{
		prefix: p.Prefix,
		player: p.Player,
		logger: p.Logger,
	})
}

type registerMessageCreateHandlersParams struct {
	prefix string
	player handlers.Player
	logger *slog.Logger
}

func registerMessageCreateHandlers(session *discordgo.Session, p registerMessageCreateHandlersParams) {
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		var err error
		switch {
		case strings.HasPrefix(m.Content, p.prefix+"ping"):
			err = handlers.Ping(s, m)
		case strings.HasPrefix(m.Content, p.prefix+"play"):
			err = middleware.MessageCreateVoiceRequired(s, m, p.player.Play)
		case strings.HasPrefix(m.Content, p.prefix+"stop"):
			err = middleware.MessageCreateVoiceRequired(s, m, p.player.Stop)
		case strings.HasPrefix(m.Content, p.prefix+"skip"):
			err = middleware.MessageCreateVoiceRequired(s, m, p.player.Skip)
		case strings.HasPrefix(m.Content, p.prefix+"queue"):
			err = middleware.MessageCreateVoiceRequired(s, m, p.player.Queue)
		}

		if err != nil {
			p.logger.Error("handle command error",
				slog.String("command", m.Content),
				log.SlogError(err),
			)
		}
	})
}
