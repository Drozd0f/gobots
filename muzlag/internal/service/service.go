package service

import (
	"errors"
	"log/slog"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

var (
	ErrGuildQueueNotFound = errors.New("guild queue not found")
)

type Service struct {
	logger *slog.Logger

	dl     ytdl.DL
	ffmpeg ffmpeg.Ffmpeg

	queue queue.Queue
}

func NewService(logger *slog.Logger, dl ytdl.DL, f ffmpeg.Ffmpeg) *Service {
	return &Service{
		logger: logger,
		dl:     dl,
		ffmpeg: f,
		queue:  queue.NewQueue(),
	}
}
