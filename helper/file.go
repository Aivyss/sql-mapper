package helper

import (
	"github.com/aivyss/sql-mapper/errors"
	"os"
	"path/filepath"
)

func ReadFile(filePath string) ([]byte, *string, errors.Error) {
	bs, err := os.ReadFile(filePath)

	if err != nil {
		return nil, nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, errors.BuildBasicErr(errors.FileReadErr)
	}
	return bs, &absFilePath, nil
}
