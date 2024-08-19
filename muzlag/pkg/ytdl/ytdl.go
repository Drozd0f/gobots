package ytdl

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ytdl struct {
	alias string
}

func NewDL(cfg Config) DL {
	dl := &ytdl{alias: cfg.Alias}

	return dl
}

func (dl *ytdl) GetAudioStream(url string, opts ...Option) (Stream, error) {
	command := []string{dl.alias}

	extractedURL, err := ExtractURL(url)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(&command)
	}

	command = append(command, extractedURL)

	return &stream{cmd: exec.Command(dl.alias, command...)}, nil
}

func (dl *ytdl) GetVideoAttributes(url string, opts ...Option) (VideoAttributes, error) {
	extractedURL, err := ExtractURL(url)
	if err != nil {
		return VideoAttributes{}, err
	}

	command := make([]string, 0, len(opts)+1)
	command = append(command, "--skip-download")
	for _, opt := range opts {
		opt(&command)
	}

	command = append(command, extractedURL)
	cmd := exec.Command(dl.alias, command...)
	out, err := cmd.Output()
	if err != nil {
		return VideoAttributes{}, fmt.Errorf("combined output: %w", err)
	}

	var va VideoAttributes
	if err = json.Unmarshal(out, &va); err != nil {
		return VideoAttributes{}, fmt.Errorf("json unmarshal: %w", err)
	}

	return va, nil
}
