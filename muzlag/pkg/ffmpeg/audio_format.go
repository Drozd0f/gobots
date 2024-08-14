package ffmpeg

type AudioFormat string

const (
	S16LeAudioFormat AudioFormat = "s16le"
)

func (f AudioFormat) String() string {
	return string(f)
}
