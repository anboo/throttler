package storage

import (
	"context"
	"sync"
	"time"

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

func (i *InMemory) ReserveRequestForQueue(ctx context.Context, workerId string, limit int) ([]Request, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var (
		res []Request
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

func (i *InMemory) UpdateStatus(ctx context.Context, id string, status Status) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	for j, r := range i.requests {
		if r.ID == id {
			i.requests[j].Status = status
			return nil
		}
	}

	return nil
}

func (i *InMemory) RequeueIdleRequests(ctx context.Context, interval time.Duration) error {
	return nil
}

func (i *InMemory) RunQueueHealthCheck(ctx context.Context, workerId string) error {
	return nil
}
