package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/anboo/throttler/service"
	database_rate "github.com/anboo/throttler/service/rate"
	"golang.org/x/time/rate"
)

func (r *Resources) initRateLimiter() error {
	switch r.Env.RateLimitStrategy {
	case "realtime":
		r.RateLimiter = service.NewRealtimeRateLimiter(r.Env.RequestsLimit, r.Env.RequestsLimitPerInterval)
		break
	case "linear":
		r.RateLimiter = rate.NewLimiter(
			rate.Every(r.Env.RequestsLimitPerInterval/time.Duration(r.Env.RequestsLimit)),
			r.Env.RequestsLimit,
		)
		break
	case "database":
		if r.Db == nil {
			log.Fatal("for database limiter need use postgres storage")
		}
		r.RateLimiter = database_rate.NewDatabaseRateLimiter(
			r.Db,
			r.Env.RequestsLimit,
			r.Env.RequestsLimitPerInterval,
		)
		break
	default:
		return fmt.Errorf("unexpected RATE_LIMIT_STRATEGY %s", r.Env.RateLimitStrategy)
	}

	return nil
}
