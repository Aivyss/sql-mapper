package context

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type TxBlock func(ctx context.Context, tx *sqlx.Tx) error
type TxBlockAuto func(ctx context.Context) error
type TxConsumer func(tx *sqlx.Tx) error

type TxManager interface {
	Txx(ctx context.Context, txBlock TxBlockAuto) error
	TxxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlockAuto) error
	Tx(ctx context.Context, txBlock TxBlock) error
	TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) error
}

func NewTxManager(writeDB *sqlx.DB) TxManager {
	return &defaultTxManager{
		writeDB: writeDB,
	}
}
