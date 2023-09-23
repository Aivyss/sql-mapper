package context

import (
	"github.com/jmoiron/sqlx"
	"sql-mapper/entity"
)

func getQueryKey(identifier string, readOnly bool) queryClientKey {
	return queryClientKey{
		identifier: identifier,
		readOnly:   readOnly,
	}
}

type dBManager int

const dbM dBManager = 1

func (ctx *dBManager) GetDB(readDB bool) *sqlx.DB {
	return GetApplicationContext().GetDB(readDB)
}

func (ctx *dBManager) GetDBs() *entity.DbSet {
	return GetApplicationContext().GetDBs()
}
