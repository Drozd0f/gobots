package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func reply(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	if _, err := s.ChannelMessageSendReply(m.ChannelID, content, m.Reference()); err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	return nil
}
