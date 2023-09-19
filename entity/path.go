package entity

import (
	"fmt"
	"sort"
	"sql-mapper/enum"
)

type Path struct {
	FilePath   string
	SqlTagName string
	DmlEnum    enum.QueryEnum
	Conditions string
}

type rawPath struct {
	FilePath   string
	SqlTagName string
	DmlEnum    enum.QueryEnum
	Conditions []*Condition
}

func NewRawPath(filePath string, tagName string, dmlEnum enum.QueryEnum, conditions ...*Condition) *rawPath {
	sort.Slice(conditions, func(i, j int) bool {
		return fmt.Sprintf("%v-%v", conditions[i].PartName, conditions[i].CaseName) <
			fmt.Sprintf("%v-%v", conditions[j].PartName, conditions[j].CaseName)
	})

	return &rawPath{
		FilePath:   filePath,
		SqlTagName: tagName,
		DmlEnum:    dmlEnum,
		Conditions: conditions,
	}
}

func (p rawPath) ToPath() Path {
	conditionKey := ""
	for _, condition := range p.Conditions {
		conditionKey += fmt.Sprintf("part:%v-case:%v", condition.PartName, condition.CaseName)
	}

	return Path{
		FilePath:   p.FilePath,
		SqlTagName: p.SqlTagName,
		DmlEnum:    p.DmlEnum,
		Conditions: conditionKey,
	}
}
