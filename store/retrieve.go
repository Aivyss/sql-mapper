package store

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

func RetrieveQuery(identifier string, queryEnum enum.QueryEnum, tagName string) (*string, errors.Error) {
	queryMap := queryStore[identifier]
	var sql string

	switch queryEnum {
	case enum.SELECT:
		sql = queryMap.SelectMap[fmt.Sprintf(enum.SelectPathFormat, queryMap.FilePath, tagName)].Sql
	case enum.INSERT:
		sql = queryMap.InsertMap[fmt.Sprintf(enum.InsertPathFormat, queryMap.FilePath, tagName)].Sql
	case enum.UPDATE:
		sql = queryMap.UpdateMap[fmt.Sprintf(enum.UpdatePathFormat, queryMap.FilePath, tagName)].Sql
	case enum.DELETE:
		sql = queryMap.DeleteMap[fmt.Sprintf(enum.DeletePathFormat, queryMap.FilePath, tagName)].Sql
	case enum.CREATE, enum.DROP:
		func() {}()
	default:
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return &sql, nil
}

func RetrieveQueryMap(identifier string) (*entity.QueryMap, errors.Error) {
	queryMap, ok := queryStore[identifier]
	if !ok {
		return nil, errors.BuildBasicErr(errors.FindQueryMapErr)
	}

	return queryMap, nil
}
