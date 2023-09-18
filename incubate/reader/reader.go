package reader

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sql-mapper/errors"
	"sql-mapper/helper"
	"sql-mapper/incubate/entity"
)

func ReadMapperFile(filePath string) (*entity.DMLBody, errors.Error) {
	xmlByteSlice, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, errors.BuildBasicErr(errors.FileReadErr)
	}

	body := new(dmlBodyRaw)
	err = xml.Unmarshal(xmlByteSlice, body)
	if err != nil {
		return nil, errors.BuildErrWithOriginal(errors.ReadBodyErr, err)
	}

	for _, raw := range body.SelectRaws {
		raw.CharData = helper.ReplaceNewLineAndTabToSpace(raw.CharData)
	}
	for _, raw := range body.InputRaws {
		raw.Sql = helper.ReplaceNewLineAndTabToSpace(raw.Sql)
	}
	for _, raw := range body.DeleteRaws {
		raw.Sql = helper.ReplaceNewLineAndTabToSpace(raw.Sql)
	}
	for _, raw := range body.UpdateRaws {
		raw.Sql = helper.ReplaceNewLineAndTabToSpace(raw.Sql)
	}

	return body.toEntity(absFilePath)
}
