package entity

import (
	"fmt"
	"sql-mapper/entity"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

type QueryMap struct {
	FilePath  string
	SelectMap map[string]*Select        // key: filePath/SELECT/tagName value: *Select
	InsertMap map[string]*entity.Insert // key: filePath/INSERT/tagName value: *Insert
	UpdateMap map[string]*entity.Update // key: filePath/UPDATE/tagName value: *Update
	DeleteMap map[string]*entity.Delete // key: filePath/DELETE/tagName value: *Delete
}

func (m *QueryMap) FindQueryInSelect(tagName string) (*Select, errors.Error) {
	query := m.SelectMap[fmt.Sprintf(enum.SelectPathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInInsert(tagName string) (*entity.Insert, errors.Error) {
	query := m.InsertMap[fmt.Sprintf(enum.InsertPathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInUpdate(tagName string) (*entity.Update, errors.Error) {
	query := m.UpdateMap[fmt.Sprintf(enum.UpdatePathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInDelete(tagName string) (*entity.Delete, errors.Error) {
	query := m.DeleteMap[fmt.Sprintf(enum.DeletePathFormat, m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}
