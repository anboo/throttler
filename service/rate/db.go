package rate

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultRateLimitId = "default"

type RateLimiter struct {
	Id                 string
	Tokens             int
	LimitTokens        int
	Interval           int
	LastReservedAt     int
	LastRecalculatedAt int
}

type DatabaseRateLimiter struct {
	db       *gorm.DB
	limit    int
	interval time.Duration
}

func NewDatabaseRateLimiter(db *gorm.DB, limit int, interval time.Duration) *DatabaseRateLimiter {
	return &DatabaseRateLimiter{
		db:       db,
		interval: interval,
		limit:    limit,
	}
}

func (l *DatabaseRateLimiter) Wait(ctx context.Context) error {
	var (
		res RateLimiter
		err error
	)

	now := currentTimestamp()

	err = l.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"limit_tokens": l.limit,
			"interval":     l.interval.Milliseconds(),
		}),
	}).Create(&RateLimiter{Id: DefaultRateLimitId}).Error
	if err != nil {
		return errors.Wrap(err, "try create lock row")
	}

	tx := l.db.Begin()
	defer func() {
		r := recover()
		if r != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = tx.Raw(
		`UPDATE rate_limiters SET
			tokens = CASE WHEN last_recalculated_at > 0 THEN tokens + TRUNC((? - last_recalculated_at) / interval) ELSE tokens END,
			last_recalculated_at = ?
		WHERE id IN (
			SELECT id FROM rate_limiters WHERE id = ? FOR UPDATE
		) RETURNING *`,
		now,
		now,
		DefaultRateLimitId,
	).Scan(&res).Error

	if err != nil {
		return errors.Wrap(err, "try fetch rate limit row")
	}

	if res.Tokens >= res.LimitTokens {
		tx.Commit()
		time.Sleep(time.Duration(res.LastReservedAt-res.Interval) * time.Millisecond)
		return l.Wait(ctx)
	}

	//err = tx.Raw(`UPDATE rate_limiters SET tokens = tokens - 1, last_reserved_at = ? WHERE id = ?`, now, DefaultRateLimitId).Error
	//if err != nil {
	//	return errors.Wrap(err, "try increment free tokens")
	//}
	tx.Commit()

	return nil
}

func currentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
