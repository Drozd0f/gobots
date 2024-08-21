package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/pkg/discordgom"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if err := discordgom.Reply(s, m, "pong"); err != nil {
		return fmt.Errorf("sending message ping: %w", err)
	}

	return nil
}
