package xml

import (
	"encoding/xml"
	"sql-mapper/entity"
	"sql-mapper/errors"
	"sql-mapper/helper"
	"sql-mapper/reader/xml/component"
)

func ReadMapperFile(filePath string) (*entity.DMLBody, errors.Error) {
	xmlByteSlice, absFilePath, err := helper.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	body := new(component.DmlBodyComponent)
	xmlErr := xml.Unmarshal(xmlByteSlice, body)
	if xmlErr != nil {
		return nil, errors.BuildErrWithOriginal(errors.ReadBodyErr, err)
	}

	for _, raw := range body.Selects {
		raw.CharData = helper.ReplaceNewLineAndTabToSpace(raw.CharData)
	}
	for _, raw := range body.Inserts {
		raw.CharData = helper.ReplaceNewLineAndTabToSpace(raw.CharData)
	}
	for _, raw := range body.Deletes {
		raw.CharData = helper.ReplaceNewLineAndTabToSpace(raw.CharData)
	}
	for _, raw := range body.Updates {
		raw.CharData = helper.ReplaceNewLineAndTabToSpace(raw.CharData)
	}

	return body.ToEntity(*absFilePath)
}

func ReadSettings(filePath string) {

}
