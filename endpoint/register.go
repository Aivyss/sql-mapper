package endpoint

import (
	"sql-mapper/errors"
	"sql-mapper/reader"
	"sql-mapper/store"
)

func Register(identifier string, filePath string) errors.Error {
	queryBody, err := reader.ReadMapperFile(filePath)
	if err != nil {
		panic(err)
	}

	return store.PersistQueries(identifier, queryBody)
}
