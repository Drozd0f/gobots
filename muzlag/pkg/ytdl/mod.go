package ytdl

import (
	"io"
)

type DL interface {
	GetAudioStream(url string, opts ...Option) (Stream, error)
	GetVideoAttributes(url string, opts ...Option) (VideoAttributes, error)
}

type Stream interface {
	GetOutput() (io.ReadCloser, error)
	Start() error
	Cancel() error
	Wait() error
}
