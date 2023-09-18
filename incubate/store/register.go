package store

import (
	"sql-mapper/errors"
	"sql-mapper/incubate/entity"
	"sql-mapper/incubate/reader"
)

func Register(identifier string, filePath string) (*entity.QueryMap, errors.Error) {
	queryBody, err := reader.ReadMapperFile(filePath)
	if err != nil {
		return nil, err
	}

	return PersistQueries(identifier, queryBody)
}
