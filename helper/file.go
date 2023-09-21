package helper

import (
	"os"
	"path/filepath"
	"sql-mapper/errors"
)

func ReadFile(filePath string) ([]byte, *string, errors.Error) {
	xmlByteSlice, err := os.ReadFile(filePath)

	if err != nil {
		return nil, nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, errors.BuildBasicErr(errors.FileReadErr)
	}
	return xmlByteSlice, &absFilePath, nil
}
