package store

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/errors"
)

const (
	SelectPathFormat = "%v/SELECT/%v"
	InsertPathFormat = "%v/INSERT/%v"
	UpdatePathFormat = "%v/UPDATE/%v"
	DeletePathFormat = "%v/DELETE/%v"
	CreatePathFormat = "%v/CREATE/%v"
	DropPathFormat   = "%v/DROP/%v"
)

var queryStore = map[string]QueryMap{} // key: filePath value: *QueryMap

type QueryMap struct {
	FilePath  string
	SelectMap map[string]*entity.Select // key: filePath/SELECT/tagName value: *entity.Select
	InsertMap map[string]*entity.Insert // key: filePath/INSERT/tagName value: *entity.Insert
	UpdateMap map[string]*entity.Update // key: filePath/UPDATE/tagName value: *entity.Update
	DeleteMap map[string]*entity.Delete // key: filePath/DELETE/tagName value: *entity.Delete
	CreateMap map[string]*entity.Create // key: filePath/CREATE/tagName value: *entity.Create
	DropMap   map[string]*entity.Drop   // key: filePath/DROP/tagName value: *entity.Drop
}

func (m *QueryMap) FindQueryInSelect(tagName string) (*entity.Select, errors.Error) {
	query := m.SelectMap[fmt.Sprintf(SelectPathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInInsert(tagName string) (*entity.Insert, errors.Error) {
	query := m.InsertMap[fmt.Sprintf(InsertPathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInUpdate(tagName string) (*entity.Update, errors.Error) {
	query := m.UpdateMap[fmt.Sprintf(UpdatePathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInDelete(tagName string) (*entity.Delete, errors.Error) {
	query := m.DeleteMap[fmt.Sprintf(DeletePathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInCreate(tagName string) (*entity.Create, errors.Error) {
	query := m.CreateMap[fmt.Sprintf(CreatePathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInDrop(tagName string) (*entity.Drop, errors.Error) {
	query := m.DropMap[fmt.Sprintf(DropPathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}
