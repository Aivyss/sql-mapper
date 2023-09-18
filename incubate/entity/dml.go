package entity

import "sql-mapper/entity"

type DMLBody struct {
	AbsFilePath string
	Selects     []*Select
	Inserts     []*entity.Insert
	Deletes     []*entity.Delete
	Updates     []*entity.Update
}

type Case struct {
	CharData string
	Name     string
}

type Part struct {
	Name     string
	CharData string
	Cases    []*Case
}

type Select struct {
	Name      string
	RawSql    string
	Parts     []*Part
	SimpleSql bool
}
