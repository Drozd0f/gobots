package middleware

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/pkg/discordgom"
	"github.com/Drozd0f/gobots/muzlag/pkg/emoji"
	"github.com/Drozd0f/gobots/muzlag/pkg/markdown"
)

type MessageCreateFunc func(s *discordgo.Session, m *discordgo.MessageCreate) error

func MessageCreateVoiceRequired(s *discordgo.Session, m *discordgo.MessageCreate, f MessageCreateFunc) error {
	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if vc := s.VoiceConnections[m.GuildID]; vc != nil {
		if vc.ChannelID != vs.ChannelID {
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
