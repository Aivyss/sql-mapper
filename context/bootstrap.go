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
	if write != nil {
		integAppCtx.txManager = NewTxManager(write)
	} else {
		integAppCtx.txManager = NewTxManager(read)
	}

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
