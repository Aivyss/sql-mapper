package store

import (
	"sql-mapper/errors"
	"sql-mapper/reader"
)

func Register(identifier string, filePath string) (*QueryMap, errors.Error) {
	queryBody, err := reader.ReadMapperFile(filePath)
	if err != nil {
		return nil, err
	}

	return PersistQueries(identifier, queryBody)
}
