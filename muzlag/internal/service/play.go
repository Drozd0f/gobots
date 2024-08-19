package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"

	"github.com/Drozd0f/gobots/muzlag/internal/queue"
	"github.com/Drozd0f/gobots/muzlag/pkg/log"
	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

const (
	frameSize int = 960           // uint16 size of each audio frame
	maxBytes      = frameSize * 4 // max size of opus data
)

func (s *Service) Play(vc *discordgo.VoiceConnection) error {
	defer func() {
		if err := s.DropGuildQueue(vc.GuildID); err != nil && !errors.Is(err, queue.ErrNotFound) {
			slog.Error("drop guild queue", log.SlogError(err))
		}
	}()

	for {
		gq, err := s.queue.GetGuildQueue(vc.GuildID)
		if err != nil {
			if errors.Is(err, queue.ErrNotFound) {
				return nil
			}

			return fmt.Errorf("getting guild queue: %w", err)
		}

		if !gq.Ready {
			return nil
		}

		if err = s.play(vc, gq); err != nil {
			return err
		}

		if len(gq.Attrs) == 0 {
			return nil
		}
	}
}

func (s *Service) play(vc *discordgo.VoiceConnection, gq *queue.GuildQueue) error {
	va, err := gq.Dequeue()
	if err != nil {
		if errors.Is(err, queue.ErrEmptyQueue) {
			return nil
		}

		return fmt.Errorf("dequeue: %w", err)
	}

	stream, err := s.dl.GetAudioStream(va.WebpageURL,
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

	send := make(chan []int16, 2)
	defer close(send)

	done := make(chan bool)
	go func() {
		s.sendPCM(vc, send)
		done <- true
	}()

	for gq.Ready && !gq.Skiped {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*2)

		err := binary.Read(pout, binary.LittleEndian, &audiobuf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}

			s.logger.Error("binary read", log.SlogError(err))
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-done:
			return nil
		}
	}

	return nil
}

// sendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func (s *Service) sendPCM(vc *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	opusEncoder, err := gopus.NewEncoder(
		s.ffmpeg.GetFrameRate(),
		s.ffmpeg.GetChannels().Int(),
		gopus.Audio,
	)
	if err != nil {
		return
	}

	for recv := range pcm {
		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return
		}

		if vc.Ready == false || vc.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		vc.OpusSend <- opus
	}
}
