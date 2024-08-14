package player

import (
	"log/slog"

	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

type NewServicePlayerParams struct {
	Logger *slog.Logger

	DL     ytdl.DL
	Ffmpeg ffmpeg.Ffmpeg
}

type ServicePlayer struct {
	logger *slog.Logger

	dl     ytdl.DL
	ffmpeg ffmpeg.Ffmpeg
}

func NewServicePlayer(p NewServicePlayerParams) *ServicePlayer {
	return &ServicePlayer{
		logger: p.Logger,
		dl:     p.DL,
		ffmpeg: p.Ffmpeg,
	}
}
