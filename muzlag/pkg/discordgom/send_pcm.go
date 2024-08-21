package discordgom

import (
	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

type SendPCMParams struct {
	Logger Logger

	FrameRate int
	Channels  int
	FrameSize int
	PCM       <-chan []int16
}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(vc *discordgo.VoiceConnection, p SendPCMParams) {
	if p.Logger == nil {
		p.Logger = noopLogger{}
	}

	if p.PCM == nil {
		return
	}

	opusEncoder, err := gopus.NewEncoder(p.FrameRate, p.Channels, gopus.Audio)
	if err != nil {
		return
	}

	maxBytes := p.FrameSize * 4
	for recv := range p.PCM {
		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, p.FrameSize, maxBytes)
		if err != nil {
			return
		}

		if vc.Ready == false || vc.OpusSend == nil {
			// Sending errors here might not be suited
			p.Logger.Warn("Discordgo not ready for opus packets",
				"ready", vc.Ready,
				"opus", opus,
			)

			return
		}

		// send encoded opus data to the sendOpus channel
		vc.OpusSend <- opus
	}
}
