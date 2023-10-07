package context

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type defaultTxManager struct {
	writeDB *sqlx.DB
}

func (d *defaultTxManager) Tx(ctx context.Context, txBlock TxBlock) (err error) {
	return d.TxWithOpt(ctx, nil, txBlock)
}

func (d *defaultTxManager) TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) error {
	tx, err := d.writeDB.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		rec := recover()
		if rec != nil {
			err2, ok := rec.(error)

			if ok {
				err = err2
			} else {
				err = errors.New("unknown err")
			}

			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = txBlock(ctx, tx)

	return err
}
