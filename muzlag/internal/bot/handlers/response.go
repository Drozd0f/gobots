package handlers

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
)

func messageSend(s *discordgo.Session, mc *discordgo.MessageCreate, content string) error {
	m, err := s.ChannelMessageSend(mc.ChannelID, content)
	if err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	go func() {
		if err = deleteBotMessage(s, m, 5*time.Second); err != nil {
			slog.Warn("failed to delete bot message: %s", log.SlogError(err))
		}
	}()

	return nil
}

func reply(s *discordgo.Session, mc *discordgo.MessageCreate, content string) error {
	m, err := s.ChannelMessageSendReply(mc.ChannelID, content, mc.Reference())
	if err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	go func() {
		if err = deleteBotMessage(s, m, 5*time.Second); err != nil {
			slog.Warn("failed to delete bot message: %s", log.SlogError(err))
		}
	}()

	return nil
}

func replyQueue(s *discordgo.Session, mc *discordgo.MessageCreate, gq *queue.GuildQueue) error {
	content := fmt.Sprintf("Now playing : %s %s", gq.CurrentAttr.Title, gq.CurrentAttr.DurationToString())
	for idx, attr := range gq.Attrs {
		content += fmt.Sprintf("\n %d - %s %s", idx+1, attr.Title, attr.DurationToString())
	}

	m, err := s.ChannelMessageSendReply(mc.ChannelID, content, mc.Reference())
	if err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	go func() {
		if err = deleteBotMessage(s, m, 30*time.Second); err != nil {
			slog.Warn("failed to delete bot message: %s", log.SlogError(err))
		}
	}()

	return nil
}

func deleteBotMessage(s *discordgo.Session, m *discordgo.Message, deleteAfter time.Duration) error {
	time.Sleep(deleteAfter)

	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		return fmt.Errorf("channel message delete: %w", err)
	}

	return nil
}
