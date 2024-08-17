package queue

import "github.com/Drozd0f/gobots/muzlag/pkg/ytdl"

type Queue interface {
	Push(id string, attr ytdl.VideoAttributes) (*GuildQueue, error)
	Drop(id string)
	GetGuildQueue(id string) (*GuildQueue, error)
}
