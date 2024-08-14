package bot

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/bot"
	"github.com/Drozd0f/gobots/muzlag/internal/bot/handlers"
	"github.com/Drozd0f/gobots/muzlag/internal/config"
	"github.com/Drozd0f/gobots/muzlag/internal/service/player"
	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	pkgLog "github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

func NewBot() *discordgo.Session {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := pkgLog.NewLogger(cfg.LogLevel)

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("error creating discord session,", err)

		return nil
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

	ph := handlers.NewPlayer(
		logger,
		player.NewServicePlayer(player.NewServicePlayerParams{
			Logger: logger,
			DL:     ytdl.NewDL(cfg.DL, ytdl.WithStandardOutput(), ytdl.WithVerbose()),
			Ffmpeg: ffmpeg.NewFfmpeg(cfg.Ffmpeg),
		}))

	bot.NewRouter(session, bot.NewRouterParams{
		Logger: logger,
		Prefix: cfg.Prefix,
		Player: ph,
	})

	if err = session.Open(); err != nil {
		log.Fatalf("session opening connection: %s", err)

		return nil
	}

	logger.Info("Bot is now running.  Press CTRL-C to exit.",
		slog.String("name", config.AppName),
	)

	return session
}
