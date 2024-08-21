package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

type PlayParams struct {
	GuildQueue      *queue.GuildQueue
	VideoAttributes ytdl.VideoAttributes

	FrameSize int // uint16 size of each audio frame
	Done      <-chan bool
	Send      chan<- []int16
}

func (s *Service) Play(p PlayParams) error {
	stream, err := s.dl.GetAudioStream(p.VideoAttributes.WebpageURL,
		ytdl.WithFormat(ytdl.BestAudioFormat),
		ytdl.WithStandardOutput(),
	)
	if err != nil {
		return fmt.Errorf("ytdl get audio stream: %w", err)
	}

	defer func() {
		if err = stream.Cancel(); err != nil {
			s.logger.Error("stream cancel", log.SlogError(err))
		}
	}()

	sout, err := stream.GetOutput()
	if err != nil {
		return fmt.Errorf("ytdl stdout pipe: %w", err)
	}

	defer func() {
		if err = sout.Close(); err != nil {
			s.logger.Error("stream out close", log.SlogError(err))
		}
	}()

	player := s.ffmpeg.PlayerFromInput(sout)
	defer func() {
		if err = player.Cancel(); err != nil {
			s.logger.Error("player close", log.SlogError(err))
		}
	}()

	pout, err := player.GetOutput()
	if err != nil {
		return fmt.Errorf("player get output: %w", err)
	}

	defer func() {
		if err = pout.Close(); err != nil {
			s.logger.Error("player output close", log.SlogError(err))
		}
	}()

	if err = stream.Start(); err != nil {
		return fmt.Errorf("stream start: %w", err)
	}

	// Starts the ffmpeg command
	if err = player.Start(); err != nil {
		return fmt.Errorf("player start: %w", err)
	}

	for p.GuildQueue.Ready && !p.GuildQueue.Skiped {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, p.FrameSize*2)

		err := binary.Read(pout, binary.LittleEndian, &audiobuf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}

			s.logger.Error("binary read", log.SlogError(err))
		}

		// Send received PCM to the sendPCM channel
		select {
		case p.Send <- audiobuf:
		case <-p.Done:
			return nil
		}
	}

	return nil
}
