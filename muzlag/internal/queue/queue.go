package queue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Drozd0f/gobots/muzlag/pkg/ytdl"
)

var (
	ErrNotFound = errors.New("not found")
)

type queue struct {
	mu *sync.Mutex

	queues map[string]*GuildQueue
}

func NewQueue() Queue {
	return &queue{
		mu:     &sync.Mutex{},
		queues: make(map[string]*GuildQueue),
	}
}

func (q *queue) Push(id string, attr ytdl.VideoAttributes) (*GuildQueue, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	gq, err := q.GetGuildQueue(id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("get guild queue: %w", err)
		}

		gq = NewGuildQueue()
		q.queues[id] = gq
	}

	gq.Enqueue(attr)

	return gq, nil
}

func (q *queue) Drop(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.queues, id)
}

func (q *queue) GetGuildQueue(id string) (*GuildQueue, error) {
	gq, exist := q.queues[id]
	if !exist {
		return nil, ErrNotFound
	}

	return gq, nil
}
