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
	var queries []entity.QueryEntity
	for _, v := range queryMap.SelectMap {
		queries = append(queries, v)
	}
	dynamicQueries := getDynamicQuery(queries)
	registerDynamicQuery(db, dynamicQueries, enum.SELECT, statementMap)
	queries = []entity.QueryEntity{}

	// create insert query statements
	for _, v := range queryMap.InsertMap {
		queries = append(queries, v)
	}
	dynamicQueries = getDynamicQuery(queries)
	registerDynamicQuery(db, dynamicQueries, enum.INSERT, statementMap)
	queries = []entity.QueryEntity{}

	// create update query statements
	for _, v := range queryMap.UpdateMap {
		queries = append(queries, v)
	}
	dynamicQueries = getDynamicQuery(queries)
	registerDynamicQuery(db, dynamicQueries, enum.UPDATE, statementMap)
	queries = []entity.QueryEntity{}

	// create delete query statements
	for _, v := range queryMap.DeleteMap {
		queries = append(queries, v)
	}
	dynamicQueries = getDynamicQuery(queries)
	registerDynamicQuery(db, dynamicQueries, enum.DELETE, statementMap)
	queries = []entity.QueryEntity{}

	return &defaultQueryClient{
		db:           db,
		queryMap:     queryMap,
		statementMap: statementMap,
	}, nil
}

func (c *defaultQueryClient) InsertOne(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error {
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
	tx, err := c.db.Beginx()
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

func getDynamicQuery(m []entity.QueryEntity) []*entity.DynamicQuery {
	var result [][]*entity.DynamicQuery

	for _, s := range m {
		var sqls []*entity.DynamicQuery

		if s.IsSimpleSql() {
			sqls = append(sqls, &entity.DynamicQuery{
				FilePath:    s.Path(),
				TagName:     s.Tag(),
				Key:         []*entity.Condition{},
				DmlEnum:     enum.SELECT,
				SqlPartials: []string{s.GetRawSql()},
			})
		} else {
			for _, part := range s.GetParts() {
				if len(part.Cases) == 0 {
					if len(sqls) == 0 {
						sqls = append(sqls, &entity.DynamicQuery{
							FilePath:    s.Path(),
							TagName:     s.Tag(),
							Key:         []*entity.Condition{},
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
								FilePath:    s.Path(),
								TagName:     s.Tag(),
								Key:         sql.Key,
								DmlEnum:     enum.SELECT,
								SqlPartials: sqlCopy,
							}

							newQuery.Key = append(newQuery.Key, &entity.Condition{
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

		result = append(result, sqls)
	}

	var flat []*entity.DynamicQuery
	for _, queries := range result {
		for _, query := range queries {
			flat = append(flat, query)
		}
	}

	return flat
}

func registerDynamicQuery(db *sqlx.DB, dynamicQueries []*entity.DynamicQuery, dmlEnum enum.QueryEnum, statementMap map[entity.Path]*sqlx.NamedStmt) {
	for _, sql := range dynamicQueries {
		rawQuery := ""
		for _, partial := range sql.SqlPartials {
			rawQuery += fmt.Sprintf("%v\n", partial)
		}

		statement, sqlxErr := db.PrepareNamed(rawQuery)
		if sqlxErr != nil {
			panic(sqlxErr)
		}

		statementMap[entity.NewRawPath(sql.FilePath, sql.TagName, dmlEnum, sql.Key...).ToPath()] = statement
	}
}

func getConditionFromPredicates(conditions []entity.PredicateConditions) []*entity.Condition {
	var cSlice []*entity.Condition
	for _, condition := range conditions {
		cc := condition()
		cSlice = append(cSlice, cc...)
	}
	return cSlice
}
