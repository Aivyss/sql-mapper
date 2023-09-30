package xml

import (
	"encoding/xml"
	"fmt"
	"github.com/aivyss/sql-mapper/entity"
	"github.com/aivyss/sql-mapper/enum"
	"github.com/aivyss/sql-mapper/errors"
	"github.com/aivyss/sql-mapper/helper"
	"github.com/aivyss/sql-mapper/reader/xml/component"
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

func ReadSettings(filePath *string) (*component.AppCtxComponent, errors.Error) {
	xmlByteSlice, absFilePath, err := helper.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	appCtxComp := new(component.AppCtxComponent)
	xmlErr := xml.Unmarshal(xmlByteSlice, appCtxComp)
	if xmlErr != nil {
		return nil, errors.BuildErrWithOriginal(errors.FileReadErr, err)
	}

	filePath = absFilePath
	return appCtxComp, nil
}

func ReadQueryMapByXml(filePath string) (*entity.QueryMap, errors.Error) {
	queryBody, err := ReadMapperFile(filePath)
	if err != nil {
		return nil, err
	}

	path := queryBody.AbsFilePath
	selectMap := map[string]*entity.Select{}
	insertMap := map[string]*entity.Insert{}
	updateMap := map[string]*entity.Update{}
	deleteMap := map[string]*entity.Delete{}

	for _, query := range queryBody.Selects {
		selectMap[fmt.Sprintf(enum.PathFormatGen.SelectPathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Inserts {
		insertMap[fmt.Sprintf(enum.PathFormatGen.InsertPathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Updates {
		updateMap[fmt.Sprintf(enum.PathFormatGen.UpdatePathFormat(), path, query.Name)] = query
	}
	for _, query := range queryBody.Deletes {
		deleteMap[fmt.Sprintf(enum.PathFormatGen.DeletePathFormat(), path, query.Name)] = query
	}

	queryMapPointer := &entity.QueryMap{
		FilePath:  path,
		SelectMap: selectMap,
		InsertMap: insertMap,
		UpdateMap: updateMap,
		DeleteMap: deleteMap,
	}

	return queryMapPointer, nil
}
