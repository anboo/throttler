package resource

import (
	"log"
	"runtime"

	"github.com/anboo/throttler/service/rate_limiter"
	"github.com/anboo/throttler/service/storage"
	"gorm.io/gorm"
)

type Resources struct {
	Db          *gorm.DB
	Storage     storage.Storage
	RateLimiter rate_limiter.RateLimiter
	Env         *ENV
}

func NewResources() *Resources {
	return &Resources{
		Env: &ENV{
			WorkersSize: runtime.GOMAXPROCS(0),
		},
	}
}

func (r *Resources) Initialize() {
	err := r.initEnv()
	if err != nil {
		log.Fatalln(err)
	}

	err = r.initStorage()
	if err != nil {
		log.Fatalln(err)
	}

	err = r.initRateLimiter()
	if err != nil {
		log.Fatalln(err)
	}
}
