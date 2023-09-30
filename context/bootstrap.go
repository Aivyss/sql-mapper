package context

import (
	"github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
	"sync"
)

var bootstrap sync.Once

func BootstrapDual(write *sqlx.DB, read *sqlx.DB) *initiator {
	ctx := GetApplicationContext()
	dbSet := ctx.GetDBs()

	dbSet.Write = write
	dbSet.Read = read

	return &initiator{}
}

func Bootstrap(db *sqlx.DB) *initiator {
	var init *initiator

	bootstrap.Do(func() {
		init = BootstrapDual(db, nil)
	})

	if init == nil {
		panic(errors.BuildBasicErr(errors.BootstrapErr))
	}

	return init
}

type initiator struct{}

func (i *initiator) InitByXml(filePath string) (ApplicationContext, errors.Error) {
	return registerXmlContext(filePath)
}
