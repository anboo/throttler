package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseStorage struct {
	db *gorm.DB
}

func NewDatabaseStorage(db *gorm.DB) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

func (d DatabaseStorage) ByID(ctx context.Context, id string) (Request, error) {
	var res Request
	err := d.db.WithContext(ctx).First(&res, "id = ?", id).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return Request{}, ErrorNotFound
	case err != nil:
		return Request{}, errors.Wrap(err, "try fetch request by id")
	}
	return res, nil
}

func (d DatabaseStorage) Create(ctx context.Context, request Request) (Request, error) {
	request.ID = uuid.New().String()
	request.Status = StatusNew
	request.CreatedAt = time.Now()

	err := d.db.WithContext(ctx).Create(&request).Error
	if err != nil {
		return Request{}, errors.Wrap(err, "try create request")
	}

	return request, nil
}

func (d DatabaseStorage) ReserveRequestForQueue(ctx context.Context, limit int) ([]Request, error) {
	var res []Request

	subQuery := d.db.Model(Request{}).Clauses(
		clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		},
	).Select("id").Where("status = ?", StatusNew).Order("id").Limit(limit)

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return d.db.Model(&res).Clauses(clause.Returning{}).Where(
			"id IN (?)", subQuery,
		).Update(
			"status", StatusInProgress,
		).Error
	})

	if err != nil {
		return res, errors.Wrap(err, "try reserve requests from db")
	}

	return res, nil
}
