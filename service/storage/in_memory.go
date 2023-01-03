package storage

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemory struct {
	lock        sync.RWMutex
	requests    map[string]*Request
	reservedIDs map[string]struct{}
}

func NewInMemory() *InMemory {
	return &InMemory{
		lock:        sync.RWMutex{},
		requests:    map[string]*Request{},
		reservedIDs: map[string]struct{}{},
	}
}

func (i *InMemory) ByID(ctx context.Context, id string) (Request, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	r, ok := i.requests[id]
	if !ok {
		return Request{}, ErrorNotFound
	}

	return *r, nil
}

func (i *InMemory) Create(ctx context.Context, request Request) (Request, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	request.ID = uuid.New().String()
	i.requests[request.ID] = &request

	return request, nil
}

func (i *InMemory) ReserveRequestForQueue(ctx context.Context, limit int) ([]Request, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var (
		res = make([]Request, limit)
		j   = 0
	)

	for _, r := range i.requests {
		_, alreadyReserved := i.reservedIDs[r.ID]
		if alreadyReserved {
			continue
		}
		res = append(res, *r)
		i.reservedIDs[r.ID] = struct{}{}

		j++
		if j >= limit {
			break
		}
	}

	return res, nil
}