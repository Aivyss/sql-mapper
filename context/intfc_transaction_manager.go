package context

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type TxBlock func(ctx context.Context, tx *sqlx.Tx) error

type TxManager interface {
	Tx(ctx context.Context, txBlock TxBlock) error
	TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) error
}

func NewTxManager(writeDB *sqlx.DB) TxManager {
	return &defaultTxManager{
		writeDB: writeDB,
	}
}
