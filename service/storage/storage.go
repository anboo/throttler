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

var (
	ErrorNotFound = errors.New("not found")
)

type Request struct {
	ID        string
	Status    Status
	CreatedAt time.Time
}

type Storage interface {
	ByID(ctx context.Context, id string) (Request, error)
	Create(ctx context.Context, request Request) (Request, error)
	ReserveRequestForQueue(ctx context.Context, limit int) ([]Request, error)
}
