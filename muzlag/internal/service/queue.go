package service

import (
	"errors"
	"fmt"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

func (s *Service) PushGuildQueue(guildID string, url string) (string, error) {
	va, err := s.dl.GetVideoAttributes(url,
		ytdl.WithTemplate(ytdl.VideoAttributesTemplate),
		ytdl.WithFormat(ytdl.BestAudioFormat),
	)
	if err != nil {
		return "", fmt.Errorf("ytdl get video attributes: %w", err)
	}

	if _, err = s.queue.Push(guildID, va); err != nil {
		return "", fmt.Errorf("queue push: %w", err)
	}

	return va.Title, nil
}

func (s *Service) GetGuildQueue(guildID string) (*queue.GuildQueue, error) {
	gq, err := s.queue.GetGuildQueue(guildID)
	if err != nil {
		if errors.Is(err, queue.ErrNotFound) {
			return nil, ErrGuildQueueNotFound
		}

		return nil, fmt.Errorf("queue get: %w", err)
	}

	return gq, nil
}

func (s *Service) SkipGuildQueue(guildID string, count int64) error {
	gq, err := s.queue.GetGuildQueue(guildID)
	if err != nil {
		return fmt.Errorf("queue get: %w", err)
	}

	if err = gq.Skip(count); err != nil {
		if errors.Is(err, queue.ErrEmptyQueue) {
			return s.DropGuildQueue(guildID)
		}

		return fmt.Errorf("queue skip: %w", err)
	}

	return nil
}

func (s *Service) DropGuildQueue(guildID string) error {
	gq, err := s.queue.GetGuildQueue(guildID)
	if err != nil {
		return fmt.Errorf("queue get: %w", err)
	}

	gq.Stop()

	s.queue.Drop(guildID)

	return nil
}
