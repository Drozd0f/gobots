package ffmpeg

import "strconv"

type Channel int

const (
	MonoChannel Channel = iota + 1
	StereoChannel
)

func (c Channel) Int() int {
	return int(c)
}

func (c Channel) String() string {
	return strconv.Itoa(c.Int())
}
