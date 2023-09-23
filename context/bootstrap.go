package context

import (
	"github.com/jmoiron/sqlx"
	"sql-mapper/errors"
)

func BootstrapDual(write *sqlx.DB, read *sqlx.DB) *initiator {
	ctx := GetApplicationContext()
	dbSet := ctx.GetDBs()

	dbSet.Write = write
	dbSet.Read = read

	return &initiator{}
}

func Bootstrap(db *sqlx.DB) *initiator {
	return BootstrapDual(db, nil)
}

type initiator struct{}

func (i *initiator) InitByXml(filePath string) (ApplicationContext, errors.Error) {
	return registerXmlContext(filePath)
}

func test() {

}
