package store

import (
	"sql-mapper/incubate/entity"
)

var queryStore = map[string]*entity.QueryMap{} // key: filePath value: *QueryMap
