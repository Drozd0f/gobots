package handlers

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/service/player"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
)

type Player struct {
	logger *slog.Logger
	s      *player.ServicePlayer
}

func NewPlayer(logger *slog.Logger, s *player.ServicePlayer) Player {
	return Player{
		logger: logger,
		s:      s,
	}
}

func (p Player) Play(s *discordgo.Session, m *discordgo.MessageCreate) error {
	sliceContent := strings.Split(m.Content, " ")
	if len(sliceContent) != 2 {
		if _, err := s.ChannelMessageSendReply(m.ChannelID, "too many words, expected: 2", m.Reference()); err != nil {
			return fmt.Errorf("channel message send reply: %w", err)
		}

		return nil
	}

	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	var vc *discordgo.VoiceConnection
	vc, err = s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		return fmt.Errorf("channel voice join: %w", err)
	}

	defer func() {
		if err = vc.Disconnect(); err != nil {
			p.logger.Error("voice connection disconnect",
				log.SlogError(err),
			)
		}
	}()

	defer vc.Close()

	if err = p.s.Play(vc, sliceContent[1]); err != nil {
		return fmt.Errorf("service play: %w", err)
	}

	return nil
}
