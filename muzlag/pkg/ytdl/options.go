package ytdl

type Option func(dl *ytdl)

var (
	WithVerbose = func() Option {
		return func(dl *ytdl) {
			dl.command = append(dl.command, "-v")
		}
	}

	WithStandardOutput = func() Option {
		return func(dl *ytdl) {
			dl.command = append(dl.command, "-o", "-")
		}
	}
)
