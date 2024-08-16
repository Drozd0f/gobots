package player

import (
	"log/slog"

	"github.com/Drozd0f/gobots/muzlag/pkg/ffmpeg"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

type ServicePlayer struct {
	logger *slog.Logger

	dl     ytdl.DL
	ffmpeg ffmpeg.Ffmpeg
}

func NewServicePlayer(logger *slog.Logger, dl ytdl.DL, f ffmpeg.Ffmpeg) *ServicePlayer {
	return &ServicePlayer{
		logger: logger,
		dl:     dl,
		ffmpeg: f,
	}
}
