package rate_limiter

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultRateLimitId = "default"

type RateLimiterRow struct {
	Id                 string
	Tokens             int
	LimitTokens        int
	Interval           int
	LastReservedAt     int
	LastRecalculatedAt int
}

func (r RateLimiterRow) TableName() string {
	return "rate_limiters"
}

type DatabaseRateLimiter struct {
	db       *gorm.DB
	limit    int
	interval time.Duration
}

func NewDatabaseRateLimiter(db *gorm.DB, limit int, interval time.Duration) *DatabaseRateLimiter {
	l := &DatabaseRateLimiter{
		db:       db,
		interval: interval,
		limit:    limit,
	}
	l.setup()

	return l
}

func (l *DatabaseRateLimiter) setup() {
	err := l.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&RateLimiterRow{
		Id:          DefaultRateLimitId,
		LimitTokens: l.limit,
		Interval:    int(l.interval.Milliseconds()),
	}).Error

	if err != nil {
		log.Fatal(errors.Wrap(err, "try create lock row"))
	}
}

func (l *DatabaseRateLimiter) Wait(ctx context.Context) error {
	type Result struct {
		result bool
		wait   int
	}

	tx := l.db.Begin()

	var res Result

	err := tx.Raw(
		`SELECT result, wait FROM take_token(?) AS (result boolean, wait BIGINT)`,
		DefaultRateLimitId,
	).Scan(&res).Error

	if err != nil {
		return errors.Wrap(err, "try fetch rate limit row")
	}

	if !res.result && res.wait == -1 {
		return errors.New("database error")
	}

	switch {
	case !res.result && res.wait == -1:
		tx.Rollback()
		return errors.New("database rate limiter error")
	case !res.result && res.wait > 0:
		tx.Commit()
		time.Sleep(time.Duration(res.wait) * time.Millisecond)
		return l.Wait(ctx)
	default:
		tx.Commit()
	}

	return nil
}
