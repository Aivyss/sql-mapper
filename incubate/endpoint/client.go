package endpoint

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sql-mapper/errors"
	entity2 "sql-mapper/incubate/entity"
	enum2 "sql-mapper/incubate/enum"
	"sql-mapper/incubate/store"
)

type defaultQueryClient struct {
	queryMap     *entity2.QueryMap
	statementMap map[entity2.Path]*sqlx.NamedStmt
	db           *sqlx.DB
}

func NewQueryClient(db *sqlx.DB, identifier string, filePath string) (QueryClient, errors.Error) {
	queryMap, err := store.Register(identifier, filePath)
	statementMap := map[entity2.Path]*sqlx.NamedStmt{}
	if err != nil {
		return nil, err
	}

	for _, s := range queryMap.SelectMap {
		var sqls []*entity2.DynamicQuery

		if s.SimpleSql {
			statement, sqlxErr := db.PrepareNamed(s.RawSql)
			if sqlxErr != nil {
				panic(sqlxErr)
			}

			statementMap[entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.SELECT).ToPath()] = statement
		} else {
			for _, part := range s.Parts {
				if len(part.Cases) == 0 {
					if len(sqls) == 0 {
						sqls = append(sqls, &entity2.DynamicQuery{
							Key:         []entity2.Condition{},
							DmlEnum:     enum2.SELECT,
							SqlPartials: []string{part.CharData},
						})
					} else {
						for _, sql := range sqls {
							sql.SqlPartials = append(sql.SqlPartials, part.CharData)
						}
					}
				} else {
					var newSqls []*entity2.DynamicQuery
					for _, sql := range sqls { // DynamicQueries
						for _, c := range part.Cases { // Cases

							var sqlCopy []string
							for _, partial := range sql.SqlPartials {
								sqlCopy = append(sqlCopy, partial)
							}

							newQuery := &entity2.DynamicQuery{
								Key:         sql.Key,
								DmlEnum:     enum2.SELECT,
								SqlPartials: sqlCopy,
							}

							newQuery.Key = append(newQuery.Key, entity2.Condition{
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
			statementMap[entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.SELECT, sql.Key...).ToPath()] = statement
		}
	}
	for _, s := range queryMap.InsertMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}
		entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.INSERT).ToPath()
		statementMap[entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.INSERT).ToPath()] = statement
	}
	for _, s := range queryMap.UpdateMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.UPDATE).ToPath()
		statementMap[entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.UPDATE).ToPath()] = statement
	}
	for _, s := range queryMap.DeleteMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[entity2.NewRawPath(queryMap.FilePath, s.Name, enum2.DELETE).ToPath()] = statement
	}

	return &defaultQueryClient{
		db:           db,
		queryMap:     queryMap,
		statementMap: statementMap,
	}, nil
}

func (c *defaultQueryClient) InsertOne(ctx context.Context, tagName string, args map[string]any) errors.Error {
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.INSERT).ToPath()
	statement := c.statementMap[path]
	_, err := statement.ExecContext(ctx, args)
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetOne(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.SELECT).ToPath()
	statement := c.statementMap[path]
	err := statement.GetContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) Get(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.SELECT).ToPath()
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

	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.SELECT).ToPath()
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
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.SELECT).ToPath()
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
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.DELETE).ToPath()
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

	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.DELETE).ToPath()
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

	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.INSERT).ToPath()
	statement := c.statementMap[path]
	reformedStatement := tx.NamedStmtContext(ctx, statement)

	_, sqlxErr := reformedStatement.ExecContext(ctx, args)
	if sqlxErr != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, sqlxErr)
	}

	return nil
}

func (c *defaultQueryClient) Update(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error) {
	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.UPDATE).ToPath()
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

	path := entity2.NewRawPath(c.queryMap.FilePath, tagName, enum2.UPDATE).ToPath()
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
