package entity

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TxContext struct {
	context.Context
	*sqlx.Tx
}

func NewTxContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return &TxContext{
		Context: ctx,
		Tx:      tx,
	}
}
