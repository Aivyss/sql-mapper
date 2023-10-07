package context

import (
	"context"
	"github.com/aivyss/sql-mapper/entity"
	"github.com/aivyss/sql-mapper/enum"
	"github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
)

func NewReadOnlyQueryClient(identifier string, filePath string) (ReadOnlyQueryClient, errors.Error) {
	return newQueryClient(identifier, filePath, true)
}

func NewQueryClient(identifier string, filePath string) (QueryClient, errors.Error) {
	return newQueryClient(identifier, filePath, false)
}

type defaultQueryClient struct {
	identifier   string
	queryMap     *entity.QueryMap
	statementMap map[entity.Path]*sqlx.NamedStmt
	readOnly     bool
}

func NewDefaultQueryClient(identifier string, queryMap *entity.QueryMap, statementMap map[entity.Path]*sqlx.NamedStmt, readOnly bool) *defaultQueryClient {
	return &defaultQueryClient{
		identifier:   identifier,
		queryMap:     queryMap,
		statementMap: statementMap,
		readOnly:     readOnly,
	}
}

func (c *defaultQueryClient) InsertOne(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	// if transaction
	txContext, ok := ctx.(*entity.TxContext)
	if ok {
		return c.InsertOneTx(txContext, txContext.Tx, tagName, args, conditions...)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.INSERT, cSlice...).ToPath()
	statement := c.statementMap[path]
	_, err := statement.ExecContext(ctx, args)
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetOne(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	// if transaction
	txContext, ok := ctx.(*entity.TxContext)
	if ok {
		return c.GetOneTx(txContext, txContext.Tx, tagName, dest, args, conditions...)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT, cSlice...).ToPath()
	statement := c.statementMap[path]
	err := statement.GetContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) Get(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	// if transaction
	txContext, ok := ctx.(*entity.TxContext)
	if ok {
		return c.GetTx(txContext, txContext.Tx, tagName, dest, args, conditions...)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT, cSlice...).ToPath()
	statement := c.statementMap[path]
	err := statement.SelectContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT, cSlice...).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	sqlxErr := reformedStatement.SelectContext(ctx, dest, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) BeginTx(_ context.Context) (*sqlx.Tx, errors.Error) {
	var dbSet = GetApplicationContext().GetDBs()
	var db *sqlx.DB

	if c.readOnly {
		db = dbSet.Read
	} else {
		db = dbSet.Write
	}
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.BuildErrWithOriginal(errors.BeginTxErr, err)
	}

	return tx, nil
}

func (c *defaultQueryClient) GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT, cSlice...).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	sqlxErr := reformedStatement.GetContext(ctx, dest, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) RollbackTx(_ context.Context, tx *sqlx.Tx) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	sqlxErr := tx.Rollback()
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.CommitTxErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) CommitTx(_ context.Context, tx *sqlx.Tx) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	sqlxErr := tx.Commit()
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.CommitTxErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) Delete(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error) {
	// if transaction
	txContext, ok := ctx.(*entity.TxContext)
	if ok {
		return c.DeleteTx(txContext, txContext.Tx, tagName, args, conditions...)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.DELETE, cSlice...).ToPath()
	statement := c.statementMap[path]

	result, err := statement.ExecContext(ctx, args)
	if err != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}
	rowsNum, err := result.RowsAffected()
	if err != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return rowsNum, nil
}

func (c *defaultQueryClient) DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.DELETE, cSlice...).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	result, sqlxErr := reformedStatement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	rowsNum, sqlxErr := result.RowsAffected()
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return rowsNum, nil
}

func (c *defaultQueryClient) InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.INSERT, cSlice...).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	_, sqlxErr := reformedStatement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) Update(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error) {
	// if transaction
	txContext, ok := ctx.(*entity.TxContext)
	if ok {
		return c.UpdateTx(txContext, txContext.Tx, tagName, args, conditions...)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.UPDATE, cSlice...).ToPath()
	statement := c.statementMap[path]

	result, sqlxErr := statement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	rowsNum, sqlxErr := result.RowsAffected()
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return rowsNum, nil
}

func (c *defaultQueryClient) UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	cSlice := getConditionFromPredicates(conditions)
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.UPDATE, cSlice...).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	result, sqlxErr := reformedStatement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	rowsNum, sqlxErr := result.RowsAffected()
	if sqlxErr != nil {
		return 0, errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return rowsNum, nil
}

func (c *defaultQueryClient) Id() string {
	return c.identifier
}

func (c *defaultQueryClient) ReadOnly() bool {
	return c.readOnly
}
