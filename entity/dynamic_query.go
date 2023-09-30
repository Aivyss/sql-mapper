package entity

import "github.com/aivyss/sql-mapper/enum"

type DynamicQuery struct {
	FilePath    string
	TagName     string
	Key         []*Condition
	DmlEnum     enum.QueryEnum
	SqlPartials []string
}
