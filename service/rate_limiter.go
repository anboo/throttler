package service

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	limit    int
	reserved int
	interval time.Duration

	ticker     *time.Ticker
	lock       sync.RWMutex
	nextTicket time.Time
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	r := &RateLimiter{
		limit:    limit,
		reserved: 0,
		interval: interval,
		ticker:   time.NewTicker(interval),
		lock:     sync.RWMutex{},
	}
	r.start()

	return r
}

func (r *RateLimiter) start() {
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

func (r *RateLimiter) Take(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.reserved >= r.limit {
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.nextTicket.Sub(time.Now())):
		}
	}

	r.nextTicket = time.Now().Add(r.interval)
	r.reserved++
}
