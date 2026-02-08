package discordgom

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
)

type SendPCMParams struct {
	Logger Logger

	FrameRate int
	Channels  int
	FrameSize int
	PCM       <-chan []int16
}

// var pcm []int16 = ... // obtain your raw PCM data somewhere
//const bufferSize = 1000 // choose any buffer size you like. 1k is plenty.
//
//// Check the frame size. You don't need to do this if you trust your input.
//frameSize := len(pcm) // must be interleaved if stereo
//frameSizeMs := float32(frameSize) / channels * 1000 / sampleRate
//switch frameSizeMs {
//case 2.5, 5, 10, 20, 40, 60:
//    // Good.
//default:
//    return fmt.Errorf("Illegal frame size: %d bytes (%f ms)", frameSize, frameSizeMs)
//}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(vc *discordgo.VoiceConnection, p SendPCMParams) {
	if p.Logger == nil {
		p.Logger = noopLogger{}
	}

	if p.PCM == nil {
		return
	}

	opusEncoder, err := opus.NewEncoder(p.FrameRate, p.Channels, opus.AppAudio)
	if err != nil {
		return
	}

	for recv := range p.PCM {
		select {
		case <-vc.Dead:
			p.Logger.Warn("voice dead")
		default:
			buf := make([]byte, p.FrameSize*2)
			// try encoding pcm frame with Opus
			n, err := opusEncoder.Encode(recv, buf)
			if err != nil {
				slog.Error("Failed to encode PCM packet", slog.Any("error", err))

				return
			}

			if vc.OpusSend == nil {
				// Sending errors here might not be suited
				p.Logger.Warn("Discordgo not ready for opus packets", "opus", buf)

				return
			}

			// send encoded opus data to the sendOpus channel
			vc.OpusSend <- buf[:n]
		}
	}
}
