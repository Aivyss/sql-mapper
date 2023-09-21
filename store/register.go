package store

import (
	"sql-mapper/entity"
	"sql-mapper/errors"
	"sql-mapper/reader/xml"
)

func Register(identifier string, filePath string) (*entity.QueryMap, errors.Error) {
	queryBody, err := xml.ReadMapperFile(filePath)
	if err != nil {
		return nil, err
	}

	return PersistQueries(identifier, queryBody)
}
