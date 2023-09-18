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
	statementMap map[entity.Path]*sqlx.NamedStmt
	db           *sqlx.DB
}

func NewQueryClient(db *sqlx.DB, identifier string, filePath string) (QueryClient, errors.Error) {
	queryMap, err := store.Register(identifier, filePath)
	statementMap := map[entity.Path]*sqlx.NamedStmt{}
	if err != nil {
		return nil, err
	}

	// create select query statements
	for _, s := range queryMap.SelectMap {
		var sqls []*entity.DynamicQuery

		if s.SimpleSql {
			statement, sqlxErr := db.PrepareNamed(s.RawSql)
			if sqlxErr != nil {
				panic(sqlxErr)
			}

			statementMap[entity.NewRawPath(queryMap.FilePath, s.Name, enum.SELECT).ToPath()] = statement
		} else {
			for _, part := range s.Parts {
				if len(part.Cases) == 0 {
					if len(sqls) == 0 {
						sqls = append(sqls, &entity.DynamicQuery{
							Key:         []entity.Condition{},
							DmlEnum:     enum.SELECT,
							SqlPartials: []string{part.CharData},
						})
					} else {
						for _, sql := range sqls {
							sql.SqlPartials = append(sql.SqlPartials, part.CharData)
						}
					}
				} else {
					var newSqls []*entity.DynamicQuery
					for _, sql := range sqls { // DynamicQueries
						for _, c := range part.Cases { // Cases

							var sqlCopy []string
							for _, partial := range sql.SqlPartials {
								sqlCopy = append(sqlCopy, partial)
							}

							newQuery := &entity.DynamicQuery{
								Key:         sql.Key,
								DmlEnum:     enum.SELECT,
								SqlPartials: sqlCopy,
							}

							newQuery.Key = append(newQuery.Key, entity.Condition{
								CaseName: c.Name,
								PartName: part.Name,
							})
							newQuery.SqlPartials = append(newQuery.SqlPartials, c.CharData)

							newSqls = append(newSqls, newQuery)
						}
					}

					sqls = newSqls
				}
			}
		}

		for _, sql := range sqls {
			rawQuery := ""
			for _, partial := range sql.SqlPartials {
				rawQuery += fmt.Sprintf("%v\n", partial)
			}

			statement, sqlxErr := db.PrepareNamed(rawQuery)
			if sqlxErr != nil {
				panic(sqlxErr)
			}
			statementMap[entity.NewRawPath(queryMap.FilePath, s.Name, enum.SELECT, sql.Key...).ToPath()] = statement
		}
	}

	// create insert query statements
	for _, s := range queryMap.InsertMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}
		entity.NewRawPath(queryMap.FilePath, s.Name, enum.INSERT).ToPath()
		statementMap[entity.NewRawPath(queryMap.FilePath, s.Name, enum.INSERT).ToPath()] = statement
	}
	for _, s := range queryMap.UpdateMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		entity.NewRawPath(queryMap.FilePath, s.Name, enum.UPDATE).ToPath()
		statementMap[entity.NewRawPath(queryMap.FilePath, s.Name, enum.UPDATE).ToPath()] = statement
	}
	for _, s := range queryMap.DeleteMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[entity.NewRawPath(queryMap.FilePath, s.Name, enum.DELETE).ToPath()] = statement
	}

	return &defaultQueryClient{
		db:           db,
		queryMap:     queryMap,
		statementMap: statementMap,
	}, nil
}

func (c *defaultQueryClient) InsertOne(ctx context.Context, tagName string, args map[string]any) errors.Error {
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.INSERT).ToPath()
	statement := c.statementMap[path]
	_, err := statement.ExecContext(ctx, args)
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetOne(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT).ToPath()
	statement := c.statementMap[path]
	err := statement.GetContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) Get(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT).ToPath()
	statement := c.statementMap[path]
	err := statement.SelectContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT).ToPath()
	statement := c.statementMap[path]
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
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.SELECT).ToPath()
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

func (c *defaultQueryClient) Delete(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error) {
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.DELETE).ToPath()
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

func (c *defaultQueryClient) DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.DELETE).ToPath()
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

func (c *defaultQueryClient) InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) errors.Error {
	if tx == nil {
		return errors.BuildBasicErr(errors.NoTxErr)
	}

	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.INSERT).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	_, sqlxErr := reformedStatement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) Update(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error) {
	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.UPDATE).ToPath()
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

func (c *defaultQueryClient) UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error) {
	if tx == nil {
		return 0, errors.BuildBasicErr(errors.NoTxErr)
	}

	path := entity.NewRawPath(c.queryMap.FilePath, tagName, enum.UPDATE).ToPath()
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
