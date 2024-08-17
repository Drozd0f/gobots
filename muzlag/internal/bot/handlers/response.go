package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
)

func reply(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	if _, err := s.ChannelMessageSendReply(m.ChannelID, content, m.Reference()); err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	return nil
}

func replyQueue(s *discordgo.Session, m *discordgo.MessageCreate, gq *queue.GuildQueue) error {
	content := fmt.Sprintf("Now playing: %s %s", gq.CurrentAttr.Title, gq.CurrentAttr.DurationToString())
	for idx, attr := range gq.Attrs {
		content += fmt.Sprintf("\n %d - %s %s", idx+1, attr.Title, attr.DurationToString())
	}

	if _, err := s.ChannelMessageSendReply(m.ChannelID, content, m.Reference()); err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	return nil
}
