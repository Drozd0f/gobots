package markdown

import (
	"fmt"
	"strings"

	"github.com/Drozd0f/gobots/muzlag/pkg/emoji"
)

var (
	Italic = func(s string) string {
		return fmt.Sprintf("*%s*", s)
	}
	Bold = func(s string) string {
		return fmt.Sprintf("**%s**", s)
	}
	Strike = func(s string) string {
		return fmt.Sprintf("~~%s~~", s)
	}
	Underline = func(s string) string {
		return fmt.Sprintf("__%s__", s)
	}
	Spoiler = func(s string) string {
		return fmt.Sprintf("||%s||", s)
	}
	Quote = func(s string) string {
		return fmt.Sprintf(">%s", s)
	}
	CodeBlock = func(s string) string {
		return fmt.Sprintf("`%s`", s)
	}
	MultiLineCodeBlock = func(s string) string {
		return fmt.Sprintf("```%s```", s)
	}
	ColoredMultiLineCodeBlock = func(style CodeBlockStyle, s string) string {
		return fmt.Sprintf("```%s\n%s```", style, s)
	}
	WithoutItalic = func(s string) string {
		return strings.ReplaceAll(s, "*", `\*`)
	}
	WithoutQuote = func(s string) string {
		return strings.ReplaceAll(s, "*", `\>`)
	}
	WithoutSpoiler = func(s string) string {
		return strings.ReplaceAll(s, "||", `\|\|`)
	}
	WithoutEmbed = func(s string) string {
		return fmt.Sprintf("<%s>", s)
	}
	GhostPing = func(s string) string {
		return fmt.Sprintf("<@&%s>", s)
	}
	QuickReaction = func(s string, e emoji.Emoji) string {
		return fmt.Sprintf("%s +%s", s, e)
	}
)
