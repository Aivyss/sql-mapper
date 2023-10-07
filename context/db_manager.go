package context

import (
	"github.com/aivyss/sql-mapper/entity"
	"github.com/jmoiron/sqlx"
)

type dBManager int

// DBManager Developers get raw DB struct
const DBManager dBManager = 1

func (ctx *dBManager) GetDB(readDB bool) *sqlx.DB {
	return GetApplicationContext().GetDB(readDB)
}

func (ctx *dBManager) GetDBs() *entity.DbSet {
	return GetApplicationContext().GetDBs()
}
