package ytdl

// Format see more format https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#format-selection
type Format string

const (
	// BestVideoFormat select the best quality format containing the video
	BestVideoFormat Format = "bestvideo"

	// BestAudioFormat select the best quality format containing audio
	BestAudioFormat Format = "bestaudio"
)

func (f Format) String() string {
	return string(f)
}
