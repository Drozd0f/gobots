package queue

import (
	"errors"
	"slices"
	"sync"

	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

var (
	ErrEmptyQueue = errors.New("queue is empty")
	//ErrQueueIsRepeated = errors.New("queue is repeated")
)

type GuildQueue struct {
	mu *sync.Mutex

	Attrs       []ytdl.VideoAttributes
	CurrentAttr ytdl.VideoAttributes

	Skiped bool
	Ready  bool
}

func NewGuildQueue(attrs ...ytdl.VideoAttributes) *GuildQueue {
	return &GuildQueue{
		mu:    &sync.Mutex{},
		Attrs: attrs,
		Ready: true,
	}
}

func (q *GuildQueue) Enqueue(attr ytdl.VideoAttributes) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Attrs = append(q.Attrs, attr)
}

func (q *GuildQueue) Dequeue() (ytdl.VideoAttributes, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Attrs) == 0 {
		return ytdl.VideoAttributes{}, ErrEmptyQueue
	}

	q.CurrentAttr, q.Attrs = q.Attrs[0], q.Attrs[1:]
	if q.Skiped {
		q.Skiped = false
	}

	return q.CurrentAttr, nil
}

func (q *GuildQueue) Skip(count int64) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	//if q.isRepeated {
	//	return ErrQueueIsRepeated
	//}

	if count != 1 {
		slices.Delete(q.Attrs, 0, int(count))
	}

	if len(q.Attrs) == 0 {
		return ErrEmptyQueue
	}

	q.Skiped = true

	return nil
}

func (q *GuildQueue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Ready = false
}
