package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusError      Status = "error"
	StatusSuccess    Status = "success"
)

var CompletedStatuses = []Status{StatusError, StatusSuccess}

var (
	ErrorNotFound = errors.New("not found")
)

type Request struct {
	ID        string
	Status    Status
	CreatedAt time.Time
}

type Worker struct {
	ID         string
	LastPingAt int
}

type Storage interface {
	ByID(ctx context.Context, id string) (Request, error)
	Create(ctx context.Context, request Request) (Request, error)
	ReserveRequestForQueue(ctx context.Context, workerId string, limit int) ([]Request, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
	RequeueIdleRequests(ctx context.Context, interval time.Duration) error
	RunQueueHealthCheck(ctx context.Context, workerId string) error
}
