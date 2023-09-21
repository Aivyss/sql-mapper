package context

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sql-mapper/endpoint"
	"sql-mapper/entity"
	"sql-mapper/enum"
	"sql-mapper/errors"
	"sql-mapper/reader/xml"
)

func NewQueryClient(db *sqlx.DB, identifier string, filePath string) (endpoint.QueryClient, errors.Error) {
	queryMap, err := xml.ReadQueryMapByXml(filePath)
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

	return endpoint.NewDefaultQueryClient(identifier, db, queryMap, statementMap), nil
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
