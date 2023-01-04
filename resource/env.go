package resource

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type ENV struct {
	DatabaseDSN                 string        `env:"DB_DSN"`
	RateLimitStrategy           string        `env:"RATE_LIMIT_STRATEGY" envDefault:"realtime"`
	RequestsLimitPerInterval    time.Duration `env:"RATE_LIMIT_REQUESTS_LIMIT_PER_INTERVAL,required" envDefault:"1m"`
	RequestsLimit               int           `env:"RATE_LIMIT_REQUESTS_INTERVAL,required" envDefault:"100"`
	IntervalCheckingNewRequests time.Duration `env:"QUEUE_INTERVAL_CHECKING_NEW_REQUESTS,required" envDefault:"1s"`
	WorkersSize                 int           `env:"QUEUE_WORKERS_SIZE"`
	HealthCheckInterval         time.Duration `env:"QUEUE_HEALTH_CHECK_INTERVAL" envDefault:"15s"`
}

func (r *Resources) initEnv() error {
	if err := env.Parse(r.Env); err != nil {
		return err
	}

	if r.Env.WorkersSize < 1 {
		log.Fatalln("worker size is not positive")
	}

	return nil
}
