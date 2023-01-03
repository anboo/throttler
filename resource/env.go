package resource

import (
	"github.com/caarlos0/env/v6"
)

func (r *Resources) initEnv() error {
	if err := env.Parse(r.Env); err != nil {
		return err
	}
	return nil
}
