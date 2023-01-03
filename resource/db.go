package resource

import (
	"strings"

	"github.com/anboo/throttler/service/storage"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (r *Resources) initStorage() error {
	if strings.HasPrefix(r.Env.DatabaseDSN, "in-memory://") {
		r.Storage = storage.NewInMemory()
	} else {
		db, err := gorm.Open(postgres.Open(r.Env.DatabaseDSN), &gorm.Config{})
		if err != nil {
			return errors.Wrap(err, "init storage gorm connect")
		}

		sqlDB, err := db.DB()
		if err != nil {
			return errors.Wrap(err, "fetch db")
		}

		err = goose.Up(sqlDB, "./migrations")
		if err != nil {
			return errors.Wrap(err, "db migration")
		}

		r.Storage = storage.NewDatabaseStorage(db)
	}
	return nil
}
