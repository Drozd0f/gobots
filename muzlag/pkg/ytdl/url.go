package ytdl

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrIvalidURLFormat = errors.New("invalid url format")

var urlre = regexp.MustCompile(`(?:https?://)?(?:www\.|m\.)?youtu(?:\.be/|be\.com/(?:watch\?(?:feature=[a-z_]+&)?v=|v/|embed/|user/(?:[^\s]+/)+|shorts/))([^?&/\s]+)`)

func ExtractURL(url string) (string, error) {
	match := urlre.FindStringSubmatch(url)

	if len(match) != 2 {
		return "", ErrIvalidURLFormat
	}

	return fmt.Sprintf("http://www.youtube.com/watch?v=%s", match[1]), nil
}
