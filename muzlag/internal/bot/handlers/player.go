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
		return reply(s, m, "too many arguments, expected: 1")
	}

	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if vc := s.VoiceConnections[m.GuildID]; vc != nil {
		if vc.ChannelID != vs.ChannelID {
			return reply(s, m, "I'm in another voice channel!")
		}

		return nil
	}

	vc, err := s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
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

func (p Player) Stop(s *discordgo.Session, m *discordgo.MessageCreate) error {
	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if s.VoiceConnections[m.GuildID] == nil {
		return reply(s, m, "I'm not in voice channel")
	}

	vc := s.VoiceConnections[m.GuildID]
	if vs.ChannelID != vc.ChannelID {
		return reply(s, m, "You can't stop me I'm in another voice channel!")
	}

	// TODO: Stop player in this voice
	return vc.Disconnect()
}
