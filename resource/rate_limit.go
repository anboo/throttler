package resource

import (
	"fmt"
	"time"

	"github.com/anboo/throttler/service"
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
	default:
		return fmt.Errorf("unexpected RATE_LIMIT_STRATEGY %s", r.Env.RateLimitStrategy)
	}

	return nil
}
