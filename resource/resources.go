package resource

import (
	"log"
	"runtime"

	"github.com/anboo/throttler/service"
	"github.com/anboo/throttler/service/storage"
)

type Resources struct {
	Storage     storage.Storage
	RateLimiter service.RateLimiter
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
