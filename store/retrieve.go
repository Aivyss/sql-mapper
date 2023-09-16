package store

import (
	"fmt"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

func RetrieveQuery(identifier string, queryEnum enum.QueryEnum, tagName string) (*string, errors.Error) {
	queryMap := queryStore[identifier]
	var sql string

	switch queryEnum {
	case enum.SELECT:
		sql = queryMap.SelectMap[fmt.Sprintf(SelectPathFormat, queryMap.FilePath, tagName)].Sql
	case enum.INSERT:
		sql = queryMap.InsertMap[fmt.Sprintf(InsertPathFormat, queryMap.FilePath, tagName)].Sql
	case enum.UPDATE:
		sql = queryMap.UpdateMap[fmt.Sprintf(UpdatePathFormat, queryMap.FilePath, tagName)].Sql
	case enum.DELETE:
		sql = queryMap.DeleteMap[fmt.Sprintf(DeletePathFormat, queryMap.FilePath, tagName)].Sql
	case enum.CREATE:
		sql = queryMap.CreateMap[fmt.Sprintf(CreatePathFormat, queryMap.FilePath, tagName)].Sql
	case enum.DROP:
		sql = queryMap.DropMap[fmt.Sprintf(DropPathFormat, queryMap.FilePath, tagName)].Sql
	default:
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return &sql, nil
}

func RetrieveQueryMap(identifier string) (*QueryMap, errors.Error) {
	queryMap, ok := queryStore[identifier]
	if !ok {
		return nil, errors.BuildBasicErr(errors.FindQueryMapErr)
	}

	return queryMap, nil
}
