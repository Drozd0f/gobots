package ffmpeg

import "strconv"

type Channel int

const (
	MonoChannel   Channel = 1
	StereoChannel Channel = 2
)

func (c Channel) Int() int {
	return int(c)
}

func (c Channel) String() string {
	return strconv.Itoa(c.Int())
}
