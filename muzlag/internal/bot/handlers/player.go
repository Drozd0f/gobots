package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/service"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/stringm"
)

type Player struct {
	logger *slog.Logger
	s      *service.Service
}

func NewPlayer(logger *slog.Logger, s *service.Service) Player {
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

		title, err := p.s.PushGuildQueue(vc.GuildID, sliceContent[1])
		if err != nil {
			return fmt.Errorf("push guild queue: %w", err)
		}

		return reply(s, m, fmt.Sprintf("Song with title: %s, add in queue", title))
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

	title, err := p.s.PushGuildQueue(vc.GuildID, sliceContent[1])
	if err != nil {
		return fmt.Errorf("push guild queue: %w", err)
	}

	if err = reply(s, m, fmt.Sprintf("Song with title: %s, add in queue", title)); err != nil {
		return err
	}

	// Send "speaking" packet over the voice websocket
	if err = vc.Speaking(true); err != nil {
		return fmt.Errorf("voice connection speaking: %w", err)
	}

	// Send not "speaking" packet over the websocket when finish
	defer func() {
		if err = vc.Speaking(false); err != nil {
			p.logger.Error("voice connection speaking false", slog.Any("error", err))
		}
	}()

	if err = p.s.Play(vc); err != nil {
		return fmt.Errorf("service play: %w", err)
	}

	if err = vc.Disconnect(); err != nil {
		return fmt.Errorf("voice connection disconnect: %w", err)
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

	if err = p.s.DropGuildQueue(vc.GuildID); err != nil {
		return fmt.Errorf("drop guild queue: %w", err)
	}

	return nil
}

func (p Player) Skip(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	sliceContent := strings.Split(m.Content, " ")
	if len(sliceContent) > 2 {
		return reply(s, m, "too many arguments, expected: 1")
	}

	var count int64 = 1
	if len(sliceContent) == 2 {
		count, err = stringm.ToInt64(sliceContent[1])
		if err != nil {
			return fmt.Errorf("stringm to int64: %w", err)
		}
	}

	if err = p.s.SkipGuildQueue(vc.GuildID, count); err != nil {
		return fmt.Errorf("skip guild queue: %w", err)
	}

	return nil
}

func (p Player) Queue(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	sliceContent := strings.Split(m.Content, " ")
	if len(sliceContent) > 2 {
		return reply(s, m, "too many arguments, expected: 1")
	}

	gq, err := p.s.GetGuildQueue(vc.GuildID)
	if err != nil {
		if errors.Is(err, service.ErrGuildQueueNotFound) {
			return reply(s, m, "I'm not playing now")
		}

		return fmt.Errorf("get attributes: %w", err)
	}

	return replyQueue(s, m, gq)
}
