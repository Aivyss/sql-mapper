package endpoint

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sql-mapper/enum"
	"sql-mapper/errors"
	"sql-mapper/store"
)

type defaultQueryClient struct {
	queryMap     *store.QueryMap
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
	for fullPath, s := range queryMap.CreateMap {
		statement, sqlxErr := db.PrepareNamed(s.Sql)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[fullPath] = statement
	}
	for fullPath, s := range queryMap.DropMap {
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

func (c *defaultQueryClient) GetOneByTagName(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	statement := c.statementMap[fmt.Sprintf(store.SelectPathFormat, c.queryMap.FilePath, tagName)]

	err := statement.GetContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetByTagName(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error {
	statement := c.statementMap[fmt.Sprintf(store.SelectPathFormat, c.queryMap.FilePath, tagName)]

	err := statement.SelectContext(ctx, dest, args) // execute
	if err != nil {
		return errors.BuildErrWithOriginal(errors.ExecuteQueryErr, err)
	}

	return nil
}

func (c *defaultQueryClient) GetRawQuery(tagName string, e enum.QueryEnum) (*string, errors.Error) {
	var sql string

	switch e {
	case enum.SELECT:
		sql = c.queryMap.SelectMap[fmt.Sprintf(store.SelectPathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.INSERT:
		sql = c.queryMap.InsertMap[fmt.Sprintf(store.InsertPathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.UPDATE:
		sql = c.queryMap.UpdateMap[fmt.Sprintf(store.UpdatePathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.DELETE:
		sql = c.queryMap.DeleteMap[fmt.Sprintf(store.DeletePathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.CREATE:
		sql = c.queryMap.CreateMap[fmt.Sprintf(store.CreatePathFormat, c.queryMap.FilePath, tagName)].Sql
	case enum.DROP:
		sql = c.queryMap.DropMap[fmt.Sprintf(store.DropPathFormat, c.queryMap.FilePath, tagName)].Sql
	default:
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return &sql, nil
}
