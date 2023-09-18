package entity

import "sql-mapper/enum"

type DynamicQuery struct {
	FilePath    string
	TagName     string
	Key         []Condition
	DmlEnum     enum.QueryEnum
	SqlPartials []string
}
