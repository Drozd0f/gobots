package ytdl

import (
	"io"
)

type DL interface {
	GetAudioStream(url string, f Format) (Stream, error)
}

type Stream interface {
	GetOutput() (io.ReadCloser, error)
	Start() error
	Cancel() error
	Wait() error
}
