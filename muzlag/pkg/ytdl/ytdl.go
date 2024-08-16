package ytdl

import (
	"os/exec"
)

type ytdl struct {
	alias   string
	command []string
}

func NewDL(cfg Config, opts ...Option) DL {
	dl := &ytdl{alias: cfg.Alias}
	for _, opt := range opts {
		opt(dl)
	}

	return dl
}

func (dl *ytdl) GetAudioStream(url string, f Format) (Stream, error) {
	args := make([]string, len(dl.command)+2)
	copy(args, dl.command)

	extractedURL, err := extractURL(url)
	if err != nil {
		return nil, err
	}

	args = append(args, "-f", f.String(), extractedURL)

	return &stream{cmd: exec.Command(dl.alias, args...)}, nil
}
