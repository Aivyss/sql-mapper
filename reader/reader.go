package reader

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sql-mapper/errors"
)

func ReadMapperFile(filePath string) (*Body, errors.Error) {
	xmlByteSlice, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	body := new(bodyRaw)
	err = xml.Unmarshal(xmlByteSlice, body)
	if err != nil {
		return nil, errors.BuildBasicErr(errors.ReadBodyErr)
	}

	return body.toEntity(absPath), nil
}
