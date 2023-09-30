package context

import (
	"github.com/aivyss/sql-mapper/entity"
	"github.com/jmoiron/sqlx"
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
