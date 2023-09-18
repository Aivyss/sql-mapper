package entity

import "sql-mapper/enum"

type DynamicQuery struct {
	Key         []Condition
	DmlEnum     enum.QueryEnum
	SqlPartials []string
}
