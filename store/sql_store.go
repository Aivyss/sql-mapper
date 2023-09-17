package store

import (
	"sql-mapper/entity"
)

var queryStore = map[string]*entity.QueryMap{} // key: filePath value: *QueryMap
