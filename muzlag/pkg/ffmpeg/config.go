package ffmpeg

type Config struct {
	Alias       string      `required:"true"`
	AudioFormat AudioFormat `required:"true" split_words:"true"`
	FrameRate   int         `required:"true" split_words:"true"`
	Channels    Channel     `required:"true"`
	BufferSize  int         `required:"true" split_words:"true"`
}
