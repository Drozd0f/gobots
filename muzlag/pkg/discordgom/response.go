package discordgom

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/emoji"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/markdown"
)

func MessageSend(s *discordgo.Session, mc *discordgo.MessageCreate, content string) error {
	if _, err := s.ChannelMessageSend(mc.ChannelID, content); err != nil {
		return fmt.Errorf("channel message send reply: %w", err)
	}

	return nil
}

func Reply(s *discordgo.Session, mc *discordgo.MessageCreate, content string) error {
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

func ReplyQueue(s *discordgo.Session, mc *discordgo.MessageCreate, gq *queue.GuildQueue) error {
	if len(gq.Attrs) == 0 {
		m, err := s.ChannelMessageSendReply(mc.ChannelID, fmt.Sprintf("%s %s %s %s",
			emoji.AngerDefaultEmoji,
			emoji.JapaneseGoblinDefaultEmoji,
			markdown.Bold("One song is not enough for queue"),
			emoji.AngerDefaultEmoji,
		), mc.Reference())
		if err != nil {
			return fmt.Errorf("channel message send reply: %w", err)
		}

		go func() {
			if err = deleteBotMessage(s, m, 10*time.Second); err != nil {
				slog.Warn("failed to delete bot message: %s", log.SlogError(err))
			}
		}()

		return nil
	}

	content := fmt.Sprintf("%s Current queue %s\n", s.State.User.ID, emoji.CoffeeDefaultEmoji)
	content = fmt.Sprintf("%s %s : %s %s",
		emoji.NqnmCenterEmoji,
		markdown.Bold(markdown.CodeBlock("Now playing")),
		gq.CurrentAttr.Title,
		gq.CurrentAttr.DurationToString(),
	)
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
