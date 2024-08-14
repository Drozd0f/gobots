package player

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"

	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

const (
	frameSize int = 960           // uint16 size of each audio frame
	maxBytes      = frameSize * 4 // max size of opus data
)

func (s *ServicePlayer) Play(vc *discordgo.VoiceConnection, url string) error {
	stream, err := s.dl.GetAudioStream(url, ytdl.BestAudioFormat)
	if err != nil {
		return fmt.Errorf("ytdl get audio stream: %w", err)
	}

	defer func() {
		if err = stream.Cancel(); err != nil {
			s.logger.Error("stream cancel", slog.Any("error", err))
		}
	}()

	sout, err := stream.GetOutput()
	if err != nil {
		return fmt.Errorf("ytdl stdout pipe: %w", err)
	}

	defer func() {
		if err = sout.Close(); err != nil {
			s.logger.Error("stream out close", slog.Any("error", err))
		}
	}()

	player := s.ffmpeg.PlayerFromInput(sout)
	defer func() {
		if err = player.Cancel(); err != nil {
			s.logger.Error("player close", slog.Any("error", err))
		}
	}()

	pout, err := player.GetOutput()
	if err != nil {
		return fmt.Errorf("player get output: %w", err)
	}

	defer func() {
		if err = pout.Close(); err != nil {
			s.logger.Error("player output close", slog.Any("error", err))
		}
	}()

	if err = stream.Start(); err != nil {
		return fmt.Errorf("stream start: %w", err)
	}

	// Starts the ffmpeg command
	if err = player.Start(); err != nil {
		return fmt.Errorf("player start: %w", err)
	}

	// Send "speaking" packet over the voice websocket
	if err = vc.Speaking(true); err != nil {
		return fmt.Errorf("voice connection speaking: %w", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		if err = vc.Speaking(false); err != nil {
			s.logger.Error("voice connection speaking false", slog.Any("error", err))
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	done := make(chan bool)
	go func() {
		s.sendPCM(vc, send)
		done <- true
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*2)

		err = binary.Read(pout, binary.LittleEndian, &audiobuf)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error reading from ffmpeg stdout: %w", err)
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-done:
			return nil
		}
	}
}

// sendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func (s *ServicePlayer) sendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
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

	for {
		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
