package resource

import (
	"context"
	"log"
	"time"

	"github.com/anboo/throttler/service/storage"
)

type ENV struct {
	DatabaseDSN              string        `env:"DB_DSN"`
	RequestsLimitPerInterval time.Duration `env:"REQUESTS_LIMIT_PER_INTERVAL"`
	RequestsLimit            int           `env:"REQUESTS_INTERVAL"`
}

type Resources struct {
	Storage storage.Storage
	Env     *ENV
}

func NewResources() *Resources {
	return &Resources{
		Env: &ENV{},
	}
}

func (r *Resources) Initialize(ctx context.Context) {
	err := r.initEnv()
	if err != nil {
		log.Fatalln(err)
	}

	err = r.initStorage()
	if err != nil {
		log.Fatalln(err)
	}
}
