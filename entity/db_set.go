package entity

import "github.com/jmoiron/sqlx"

type DbSet struct {
	Write *sqlx.DB
	Read  *sqlx.DB
}
