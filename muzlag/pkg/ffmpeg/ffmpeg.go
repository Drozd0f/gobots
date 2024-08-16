package ffmpeg

import (
	"io"
	"os/exec"
	"strconv"
)

type ffmpeg struct {
	alias       string
	audioFormat AudioFormat
	frameRate   int
	channels    Channel
	bufferSize  int
}

func NewFfmpeg(cfg Config) Ffmpeg {
	return &ffmpeg{
		alias:       cfg.Alias,
		audioFormat: cfg.AudioFormat,
		frameRate:   cfg.FrameRate,
		channels:    cfg.Channels,
		bufferSize:  cfg.BufferSize,
	}
}

func (f *ffmpeg) PlayerFromFile(filepath string) Player {
	return &player{
		cmd: exec.Command(f.alias,
			"-i", filepath,
			"-f", f.audioFormat.String(),
			"-ar", strconv.Itoa(f.frameRate),
			"-ac", f.channels.String(),
			"pipe:1",
		),
		bufferSize: f.bufferSize,
	}
}

func (f *ffmpeg) PlayerFromInput(r io.Reader) Player {
	p := &player{
		cmd: exec.Command(f.alias,
			"-i", "pipe:0",
			"-f", f.audioFormat.String(),
			"-ar", strconv.Itoa(f.frameRate),
			"-ac", f.channels.String(),
			"pipe:1",
		),
		bufferSize: f.bufferSize,
	}

	p.SetInput(r)

	return p
}

func (f *ffmpeg) GetChannels() Channel {
	return f.channels
}

func (f *ffmpeg) GetFrameRate() int {
	return f.frameRate
}
