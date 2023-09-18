package entity

import (
	"fmt"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

type QueryMap struct {
	FilePath  string
	SelectMap map[string]*Select // key: filePath/SELECT/tagName value: *Select
	InsertMap map[string]*Insert // key: filePath/INSERT/tagName value: *Insert
	UpdateMap map[string]*Update // key: filePath/UPDATE/tagName value: *Update
	DeleteMap map[string]*Delete // key: filePath/DELETE/tagName value: *Delete
}

func (m *QueryMap) FindQueryInSelect(tagName string) (*Select, errors.Error) {

	query := m.SelectMap[fmt.Sprintf(enum.PathFormatGen.SelectPathFormat(), m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInInsert(tagName string) (*Insert, errors.Error) {
	query := m.InsertMap[fmt.Sprintf(enum.PathFormatGen.InsertPathFormat(), m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInUpdate(tagName string) (*Update, errors.Error) {
	query := m.UpdateMap[fmt.Sprintf(enum.PathFormatGen.UpdatePathFormat(), m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}

func (m *QueryMap) FindQueryInDelete(tagName string) (*Delete, errors.Error) {
	query := m.DeleteMap[fmt.Sprintf(enum.PathFormatGen.DeletePathFormat(), m.FilePath, tagName)]
	if query == nil {
		return nil, errors.BuildBasicErr(errors.FindQueryErr)
	}

	return query, nil
}
