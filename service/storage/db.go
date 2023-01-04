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

func (d *DatabaseStorage) ByID(ctx context.Context, id string) (Request, error) {
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

func (d *DatabaseStorage) Create(ctx context.Context, request Request) (Request, error) {
	request.ID = uuid.New().String()
	request.Status = StatusNew
	request.CreatedAt = time.Now()

	err := d.db.WithContext(ctx).Create(&request).Error
	if err != nil {
		return Request{}, errors.Wrap(err, "try create request")
	}

	return request, nil
}

func (d *DatabaseStorage) ReserveRequestForQueue(ctx context.Context, workerId string, limit int) ([]Request, error) {
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
		).Updates(map[string]interface{}{
			"status":    StatusInProgress,
			"worker_id": workerId,
		}).Error
	})

	if err != nil {
		return res, errors.Wrap(err, "try reserve requests from db")
	}

	return res, nil
}

func (d *DatabaseStorage) UpdateStatus(ctx context.Context, id string, status Status) error {
	err := d.db.WithContext(ctx).Model(&Request{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return errors.Wrap(err, "update status")
	}
	return nil
}

func (d *DatabaseStorage) RequeueIdleRequests(ctx context.Context, interval time.Duration) error {
	subQuery := d.db.Model(&Worker{}).Select("id").Where("? - last_ping_at >= ? * 1.5", time.Now().Unix(), interval.Seconds())

	err := d.db.WithContext(ctx).Model(&Request{}).Where(
		"worker_id IN (?) AND status NOT IN (?)",
		subQuery,
		CompletedStatuses,
	).Update("status", StatusNew).Error

	if err != nil {
		return errors.Wrap(err, "try reset status to new")
	}

	err = d.db.Where("id IN (?)", subQuery).Delete(&Worker{}).Error
	if err != nil {
		return errors.Wrap(err, "reset idle tasks remove workers")
	}

	return nil
}

func (d *DatabaseStorage) RunQueueHealthCheck(ctx context.Context, workerId string) error {
	lastPingAt := int(time.Now().Unix())

	err := d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"last_ping_at": lastPingAt}),
	}).Create(&Worker{ID: workerId, LastPingAt: lastPingAt}).Error

	if err != nil {
		return errors.Wrap(err, "try health check scenario")
	}

	return nil
}
