package context

import (
	"context"
	"database/sql"
	"github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
)

type TxBlock func(ctx context.Context, tx *sqlx.Tx) error
type TxBlockAuto func(ctx context.Context) error

type TxManager interface {
	Txx(ctx context.Context, txBlock TxBlockAuto) errors.Error
	TxxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlockAuto) errors.Error
	Tx(ctx context.Context, txBlock TxBlock) errors.Error
	TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) errors.Error
}

func NewTxManager(writeDB *sqlx.DB) TxManager {
	return &defaultTxManager{
		writeDB: writeDB,
	}
}
