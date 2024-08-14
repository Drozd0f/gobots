package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	if err != nil {
		return fmt.Errorf("sending message ping: %w", err)
	}

	return nil
}
