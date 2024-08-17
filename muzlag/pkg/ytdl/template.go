package ytdl

import (
	"fmt"
	"strings"
	"time"
)

type Template string

func (t Template) String() string {
	return string(t)
}

const (
	VideoAttributesTemplate Template = `{
	"id": %(id)j,
	"title": %(title)j,
	"webpage_url": %(webpage_url)j,
	"duration": "%(duration)j"
}`
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	s := strings.ReplaceAll(string(b), `"`, "")

	if s == "null" || s == "NA" {
		*d = Duration(0)

		return nil
	}

	dur, err := time.ParseDuration(s + "s")
	if err != nil {
		return fmt.Errorf("time parse duration: %w", err)
	}

	*d = Duration(dur)

	return nil
}

func (d *Duration) String() string {
	return time.Duration(*d).String()
}

type VideoAttributes struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	WebpageURL string   `json:"webpage_url"`
	Duration   Duration `json:"duration"`
	StartTime  Duration `json:"start_time"`
}

func (a *VideoAttributes) DurationToString() string {
	if a.Duration == 0 {
		return "live"
	}

	return a.Duration.String()
}
