package store

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/errors"
)

func PersistQueries(identifier string, queryBody *entity.Body) (*QueryMap, errors.Error) {
	path := queryBody.AbsFilePath
	selectMap := map[string]*entity.Select{}
	insertMap := map[string]*entity.Insert{}
	updateMap := map[string]*entity.Update{}
	deleteMap := map[string]*entity.Delete{}
	createMap := map[string]*entity.Create{}
	dropMap := map[string]*entity.Drop{}

	for _, query := range queryBody.Selects {
		selectMap[fmt.Sprintf(SelectPathFormat, path, query.Name)] = &query
	}
	for _, query := range queryBody.Inserts {
		insertMap[fmt.Sprintf(InsertPathFormat, path, query.Name)] = &query
	}
	for _, query := range queryBody.Updates {
		updateMap[fmt.Sprintf(UpdatePathFormat, path, query.Name)] = &query
	}
	for _, query := range queryBody.Deletes {
		deleteMap[fmt.Sprintf(DeletePathFormat, path, query.Name)] = &query
	}
	for _, query := range queryBody.Creates {
		createMap[fmt.Sprintf(CreatePathFormat, path, query.Name)] = &query
	}
	for _, query := range queryBody.Drops {
		dropMap[fmt.Sprintf(DropPathFormat, path, query.Name)] = &query
	}

	_, ok := queryStore[identifier]
	if ok {
		return nil, errors.BuildBasicErr(errors.DuplicatedIdentifierErr)
	}

	queryMapPointer := &QueryMap{
		FilePath:  path,
		SelectMap: selectMap,
		InsertMap: insertMap,
		UpdateMap: updateMap,
		DeleteMap: deleteMap,
		CreateMap: createMap,
		DropMap:   dropMap,
	}
	queryStore[identifier] = queryMapPointer

	return queryMapPointer, nil
}
