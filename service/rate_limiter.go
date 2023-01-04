package service

import (
	"context"
	"sync"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) (err error)
}

type RealtimeRateLimiter struct {
	limit    int
	reserved int
	interval time.Duration

	ticker     *time.Ticker
	lock       sync.Mutex
	nextTicket time.Time
}

func NewRealtimeRateLimiter(limit int, interval time.Duration) *RealtimeRateLimiter {
	r := &RealtimeRateLimiter{
		limit:    limit,
		reserved: 0,
		interval: interval,
		ticker:   time.NewTicker(interval),
		lock:     sync.Mutex{},
	}
	r.start()

	return r
}

func (r *RealtimeRateLimiter) start() {
	go func() {
		for {
			r.lock.Lock()
			if r.reserved > 0 {
				r.reserved--
			}
			r.lock.Unlock()

			<-r.ticker.C
		}
	}()
}

func (r *RealtimeRateLimiter) Wait(ctx context.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.reserved >= r.limit {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(r.nextTicket.Sub(time.Now())):
		}
	}

	r.nextTicket = time.Now().Add(r.interval)
	r.reserved++

	return nil
}
