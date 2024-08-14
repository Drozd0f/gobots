package ffmpeg

import (
	"io"
)

type Ffmpeg interface {
	PlayerFromFile(filepath string) Player
	PlayerFromInput(r io.Reader) Player
	GetChannels() Channel
	GetFrameRate() int
}

type Player interface {
	SetInput(reader io.Reader)
	GetOutput() (io.ReadCloser, error)
	Start() error
	Cancel() error
	Wait() error
}
