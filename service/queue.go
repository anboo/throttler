package service

import (
	"context"
	"sync"
	"time"

	"github.com/anboo/throttler/service/storage"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Queue struct {
	checkingInterval time.Duration

	group   *errgroup.Group
	once    *sync.Once
	workers int
	queue   chan storage.Request

	storage    storage.Storage
	httpClient *HttpClient

	logger *zerolog.Logger
}

func NewQueue(
	ctx context.Context,
	checkingInterval time.Duration,
	httpClient *HttpClient,
	workers int,
	db storage.Storage,
	logger *zerolog.Logger,
) *Queue {
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(workers)

	return &Queue{
		group:            group,
		once:             &sync.Once{},
		httpClient:       httpClient,
		workers:          workers,
		checkingInterval: checkingInterval,
		queue:            make(chan storage.Request, workers),
		storage:          db,
		logger:           logger,
	}
}

func (q *Queue) reserveRequest(ctx context.Context) {
	reqs, err := q.storage.ReserveRequestForQueue(ctx, q.workers)
	if err != nil {
		q.logger.Err(err).Msg("try reserve requests for queue")
		return
	}

	for _, r := range reqs {
		select {
		case <-ctx.Done():
			err := q.storage.UpdateStatus(context.Background(), r.ID, storage.StatusNew)
			if err != nil {
				q.logger.Err(err).Msg("shutdown cancel reservation")
			}
		default:
			q.queue <- r
		}
	}
}

func (q *Queue) Start(ctx context.Context) {
	timer := time.NewTicker(q.checkingInterval)

	q.logger.Info().Str("interval", q.checkingInterval.String()).Int("workers", q.workers).Msg("start queue")

	go func() {
		for {
			q.logger.Debug().Msg("checking new jobs")
			q.reserveRequest(ctx)

			select {
			case <-timer.C:
				continue
			case <-ctx.Done():
				timer.Stop()
				q.logger.Warn().Msg("close")
				return
			}
		}
	}()

	consumerGroup, _ := errgroup.WithContext(ctx)
	consumerGroup.SetLimit(q.workers)
	for i := 0; i < q.workers; i++ {
		consumerGroup.Go(func() error {
			q.consuming(ctx)
			return nil
		})
	}
	consumerGroup.Wait()
}

func (q *Queue) consuming(ctx context.Context) {
	for {
		select {
		case req := <-q.queue:
			err := q.httpClient.Request(ctx, req)

			var status storage.Status
			if err != nil {
				status = storage.StatusError
			} else {
				status = storage.StatusSuccess
			}

			err = q.storage.UpdateStatus(ctx, req.ID, status)
			if err != nil {
				q.logger.Err(err).Str("status", string(status)).Msg("try update status")
			}
		case <-ctx.Done():
			q.logger.Warn().Msg("try graceful shutdown queue")
			q.shutdown()
			return
		}
	}
}

func (q *Queue) shutdown() {
	q.once.Do(func() {
		defer close(q.queue)
		q.group.Wait()
	})
}
