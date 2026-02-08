package middleware

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/pkg/discordgom"
	"github.com/Drozd0f/gobots/muzlag/pkg/emoji"
	"github.com/Drozd0f/gobots/muzlag/pkg/markdown"
)

type MessageCreateFunc func(s *discordgo.Session, m *discordgo.MessageCreate) error

func MessageCreateVoiceRequired(s *discordgo.Session, m *discordgo.MessageCreate, f MessageCreateFunc) error {
	botState, err := s.State.VoiceState(m.GuildID, s.State.User.ID)
	if err != nil {
		if errors.Is(err, discordgo.ErrStateNotFound) {
			return f(s, m)
		}

		return fmt.Errorf("bot get channel: %w", err)
	}

	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		if errors.Is(err, discordgo.ErrStateNotFound) {
			return discordgom.Reply(s, m,
				fmt.Sprintf("%s from where you sad that? %s",
					markdown.Bold(m.Author.Username),
					emoji.ThinkingDefaultEmoji,
				),
			)
		}

		return fmt.Errorf("get channel: %w", err)
	}

	if vc := s.VoiceConnections[m.GuildID]; vc != nil {
		if botState.ChannelID != vs.ChannelID {
			return discordgom.Reply(s, m,
				fmt.Sprintf("%s from where you sad that? %s",
					markdown.Bold(m.Author.Username),
					emoji.ThinkingDefaultEmoji,
				),
			)
		}
	}

	return f(s, m)
}
