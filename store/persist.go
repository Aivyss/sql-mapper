package store

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/errors"
)

func PersistQueries(identifier string, queryBody *entity.Body) errors.Error {
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
		return errors.BuildBasicErr(errors.RegisterQueryErr)
	}

	queryStore[identifier] = QueryMap{
		FilePath:  path,
		SelectMap: selectMap,
		InsertMap: insertMap,
		UpdateMap: updateMap,
		DeleteMap: deleteMap,
		CreateMap: createMap,
		DropMap:   dropMap,
	}

	return nil
}
