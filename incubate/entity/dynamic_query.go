package entity

import "sql-mapper/incubate/enum"

type DynamicQuery struct {
	Key         []Condition
	DmlEnum     enum.QueryEnum
	SqlPartials []string
}
