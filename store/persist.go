package store

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

func PersistQueries(identifier string, queryBody *entity.DMLBody) (*entity.QueryMap, errors.Error) {
	_, ok := queryStore[identifier]
	if ok {
		return nil, errors.BuildBasicErr(errors.DuplicatedIdentifierErr)
	}

	path := queryBody.AbsFilePath
	selectMap := map[string]*entity.Select{}
	insertMap := map[string]*entity.Insert{}
	updateMap := map[string]*entity.Update{}
	deleteMap := map[string]*entity.Delete{}

	for _, query := range queryBody.Selects {
		selectMap[fmt.Sprintf(enum.PathFormatGen.SelectPathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Inserts {
		insertMap[fmt.Sprintf(enum.PathFormatGen.InsertPathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Updates {
		updateMap[fmt.Sprintf(enum.PathFormatGen.UpdatePathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Deletes {
		deleteMap[fmt.Sprintf(enum.PathFormatGen.DeletePathFormat(), path, query.Name)] = query
	}

	queryMapPointer := &entity.QueryMap{
		FilePath:  path,
		SelectMap: selectMap,
		InsertMap: insertMap,
		UpdateMap: updateMap,
		DeleteMap: deleteMap,
	}
	queryStore[identifier] = queryMapPointer

	return queryMapPointer, nil
}
