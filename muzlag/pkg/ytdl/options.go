package ytdl

type Option func(commands *[]string)

var (
	WithStandardOutput = func() Option {
		return func(commands *[]string) {
			*commands = append(*commands, "-o", "-")
		}
	}

	WithVerbose = func() Option {
		return func(commands *[]string) {
			*commands = append(*commands, "-v")
		}
	}

	WithTemplate = func(t Template) Option {
		return func(commands *[]string) {
			*commands = append(*commands, "--print", t.String())
		}
	}

	WithFormat = func(f Format) Option {
		return func(commands *[]string) {
			*commands = append(*commands, "-f", f.String())
		}
	}
)
