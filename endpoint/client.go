package endpoint

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sql-mapper/entity"
	"sql-mapper/enum"
	"sql-mapper/errors"
	"sql-mapper/store"
)

type defaultQueryClient struct {
	queryMap     *entity.QueryMap
	statementMap map[string]*sqlx.NamedStmt
	db           *sqlx.DB
}

func NewQueryClient(db *sqlx.DB, identifier string, filePath string) (QueryClient, errors.Error) {
	queryMap, err := store.Register(identifier, filePath)
	statementMap := map[string]*sqlx.NamedStmt{}
	if err != nil {
		return nil, err
	}

	for fullPath, s := range queryMap.SelectMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[fullPath] = statement
	}
	for fullPath, s := range queryMap.InsertMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[fullPath] = statement
	}
	for fullPath, s := range queryMap.UpdateMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[fullPath] = statement
	}
	for fullPath, s := range queryMap.DeleteMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[fullPath] = statement
	}

	return &defaultQueryClient{
		db:           db,
		queryMap:     queryMap,
		statementMap: statementMap,
	}, nil
}

func (c *defaultQueryClient) InsertOne(ctx context.Context, tagName string, args map[string]any) errors.Error {
	statement := c.statementMap[fmt.Sprintf(enum.InsertPathFormat, c.queryMap.FilePath, tagName)]
	_, err := statement.ExecContext(ctx, args)
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetOne(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	statement := c.statementMap[fmt.Sprintf(enum.SelectPathFormat, c.queryMap.FilePath, tagName)]
	err := statement.GetContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) Get(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	statement := c.statementMap[fmt.Sprintf(enum.SelectPathFormat, c.queryMap.FilePath, tagName)]
	err := statement.SelectContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetRawQuery(tagName string, e enum.QueryEnum) (*string, errors.Error) {
	var query string

	switch e {
	case enum.SELECT:
		query = c.queryMap.SelectMap[fmt.Sprintf(enum.SelectPathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.INSERT:
		query = c.queryMap.InsertMap[fmt.Sprintf(enum.InsertPathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.UPDATE:
		query = c.queryMap.UpdateMap[fmt.Sprintf(enum.UpdatePathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.DELETE:
		query = c.queryMap.DeleteMap[fmt.Sprintf(enum.DeletePathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.CREATE, enum.DROP:
		func() {}()
	default:
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return &query, nil
}

func (c *defaultQueryClient) GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	statement := c.statementMap[fmt.Sprintf(enum.SelectPathFormat, c.queryMap.FilePath, tagName)]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	sqlxErr := reformedStatement.SelectContext(ctx, dest, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) BeginTx(_ context.Context) (*sqlx.Tx, errors.Error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, errors.BuildErrWithOriginal(errors.BeginTxErr, err)
	}

	return tx, nil
}

func (c *defaultQueryClient) GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	statement := c.statementMap[fmt.Sprintf(enum.SelectPathFormat, c.queryMap.FilePath, tagName)]
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

func (c *defaultQueryClient) Delete(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error) {
	statement := c.statementMap[fmt.Sprintf(enum.DeletePathFormat, c.queryMap.FilePath, tagName)]

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

func (c *defaultQueryClient) DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	statement := c.statementMap[fmt.Sprintf(enum.DeletePathFormat, c.queryMap.FilePath, tagName)]
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

func (c *defaultQueryClient) InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	rawQuery, err := c.GetRawQuery(tagName, enum.INSERT)
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	_, sqlxErr := tx.NamedExecContext(ctx, *rawQuery, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) Update(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error) {
	statement := c.statementMap[fmt.Sprintf(enum.UpdatePathFormat, c.queryMap.FilePath, tagName)]

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

func (c *defaultQueryClient) UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	statement := c.statementMap[fmt.Sprintf(enum.UpdatePathFormat, c.queryMap.FilePath, tagName)]
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
