package context

import (
	"context"
	"database/sql"
	"github.com/aivyss/sql-mapper/entity"
	lerr "github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
)

type defaultTxManager struct {
	writeDB *sqlx.DB
}

func (d *defaultTxManager) Txx(ctx context.Context, txBlock TxBlockAuto) lerr.Error {
	return d.TxxWithOpt(ctx, nil, txBlock)
}

func (d *defaultTxManager) TxxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlockAuto) (err lerr.Error) {
	tx, dbErr := d.writeDB.BeginTxx(ctx, opts)
	if dbErr != nil {
		return lerr.BuildErrWithOriginal(lerr.StartTransactionErr, err)
	}

	defer func() {
		rec := recover()
		if rec != nil {
			err2, ok := rec.(error)

			if ok {
				err = lerr.BuildErrWithOriginal(lerr.TransactionClientSidePanicErr, err2)
			} else {
				err = lerr.BuildBasicErr(lerr.UnknownTransactionErr)
			}

			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	ctx = entity.NewTxContext(ctx, tx)
	clientSideErr := txBlock(ctx)
	if clientSideErr != nil {
		err = lerr.BuildErrWithOriginal(lerr.TransactionClientSideErr, clientSideErr)
	}

	return err
}

func (d *defaultTxManager) Tx(ctx context.Context, txBlock TxBlock) (err lerr.Error) {
	return d.TxWithOpt(ctx, nil, txBlock)
}

func (d *defaultTxManager) TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) (err lerr.Error) {
	tx, dbErr := d.writeDB.BeginTxx(ctx, opts)
	if dbErr != nil {
		return lerr.BuildErrWithOriginal(lerr.StartTransactionErr, err)
	}

	defer func() {
		rec := recover()
		if rec != nil {
			err2, ok := rec.(error)

			if ok {
				err = lerr.BuildErrWithOriginal(lerr.TransactionClientSidePanicErr, err2)
			} else {
				err = lerr.BuildBasicErr(lerr.UnknownTransactionErr)
			}

			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	clientSideErr := txBlock(ctx, tx)
	if clientSideErr != nil {
		err = lerr.BuildErrWithOriginal(lerr.TransactionClientSideErr, clientSideErr)
	}

	return err
}
