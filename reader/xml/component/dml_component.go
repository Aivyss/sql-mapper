package component

import (
	"encoding/xml"
	"sql-mapper/entity"
	"sql-mapper/errors"
	"sql-mapper/helper"
)

type DmlBodyComponent struct {
	XMLName xml.Name           `xml:"Body"`
	Selects []*selectComponent `xml:"Select"`
	Inserts []*insertComponent `xml:"Insert"`
	Updates []*updateComponent `xml:"Update"`
	Deletes []*deleteComponent `xml:"Delete"`
}

func (b DmlBodyComponent) ToEntity(absFilePath string) (*entity.DMLBody, errors.Error) {
	var s []*entity.Select
	for _, sql := range b.Selects {
		elem, err := sql.toEntity(absFilePath)
		if err != nil {
			return nil, err
		}
		s = append(s, elem)
	}

	var i []*entity.Insert
	for _, sql := range b.Inserts {
		elem, err := sql.toEntity(absFilePath)
		if err != nil {
			return nil, err
		}
		i = append(i, elem)
	}

	var u []*entity.Update
	for _, sql := range b.Updates {
		elem, err := sql.toEntity(absFilePath)
		if err != nil {
			return nil, err
		}
		u = append(u, elem)
	}

	var d []*entity.Delete
	for _, sql := range b.Deletes {
		elem, err := sql.toEntity(absFilePath)
		if err != nil {
			return nil, err
		}
		d = append(d, elem)
	}

	return &entity.DMLBody{
		AbsFilePath: absFilePath,
		Selects:     s,
		Inserts:     i,
		Updates:     u,
		Deletes:     d,
	}, nil
}

type insertComponent struct {
	CharData string           `xml:",chardata"`
	Name     string           `xml:"name,attr"`
	Parts    []*partComponent `xml:"Part"`
}

func (s insertComponent) toEntity(absFilePath string) (*entity.Insert, errors.Error) {

	if helper.IsBlank(s.CharData) && len(s.Parts) == 0 {
		return nil, errors.BuildBasicErr(errors.ParseQueryErr)
	}

	part := []*entity.Part{}
	if helper.IsBlank(s.CharData) {
		for _, p := range s.Parts {
			part = append(part, p.toEntity())
		}
	}

	return &entity.Insert{
		CommonFields: &entity.CommonFields{
			FilePath:  absFilePath,
			RawSql:    s.CharData,
			SimpleSql: !helper.IsBlank(s.CharData),
			Parts:     part,
			Name:      s.Name,
		},
	}, nil
}

type updateComponent struct {
	CharData string           `xml:",chardata"`
	Name     string           `xml:"name,attr"`
	Parts    []*partComponent `xml:"Part"`
}

func (s updateComponent) toEntity(absFilePath string) (*entity.Update, errors.Error) {

	if helper.IsBlank(s.CharData) && len(s.Parts) == 0 {
		return nil, errors.BuildBasicErr(errors.ParseQueryErr)
	}

	part := []*entity.Part{}
	if helper.IsBlank(s.CharData) {
		for _, p := range s.Parts {
			part = append(part, p.toEntity())
		}
	}

	return &entity.Update{
		CommonFields: &entity.CommonFields{
			FilePath:  absFilePath,
			RawSql:    s.CharData,
			SimpleSql: !helper.IsBlank(s.CharData),
			Parts:     part,
			Name:      s.Name,
		},
	}, nil
}

type deleteComponent struct {
	CharData string           `xml:",chardata"`
	Name     string           `xml:"name,attr"`
	Parts    []*partComponent `xml:"Part"`
}

func (s deleteComponent) toEntity(absFilePath string) (*entity.Delete, errors.Error) {

	if helper.IsBlank(s.CharData) && len(s.Parts) == 0 {
		return nil, errors.BuildBasicErr(errors.ParseQueryErr)
	}

	part := []*entity.Part{}
	if helper.IsBlank(s.CharData) {
		for _, p := range s.Parts {
			part = append(part, p.toEntity())
		}
	}

	return &entity.Delete{
		CommonFields: &entity.CommonFields{
			FilePath:  absFilePath,
			RawSql:    s.CharData,
			SimpleSql: !helper.IsBlank(s.CharData),
			Parts:     part,
			Name:      s.Name,
		},
	}, nil
}

type selectComponent struct {
	XMLName  xml.Name         `xml:"Select"`
	CharData string           `xml:",chardata"`
	Name     string           `xml:"name,attr"`
	Parts    []*partComponent `xml:"Part"`
}

func (s *selectComponent) toEntity(absFilePath string) (*entity.Select, errors.Error) {
	if helper.IsBlank(s.CharData) && len(s.Parts) == 0 {
		return nil, errors.BuildBasicErr(errors.ParseQueryErr)
	}

	part := []*entity.Part{}
	if helper.IsBlank(s.CharData) {
		for _, p := range s.Parts {
			part = append(part, p.toEntity())
		}
	}

	return &entity.Select{
		CommonFields: &entity.CommonFields{
			FilePath:  absFilePath,
			RawSql:    s.CharData,
			SimpleSql: !helper.IsBlank(s.CharData),
			Parts:     part,
			Name:      s.Name,
		},
	}, nil
}

type partComponent struct {
	Name     string          `xml:"name,attr"`
	CharData string          `xml:",chardata"`
	Cases    []caseComponent `xml:"Case"`
}

func (p *partComponent) toEntity() *entity.Part {

	cases := []*entity.Case{}
	for _, c := range p.Cases {
		cases = append(cases, c.toEntity())
	}

	return &entity.Part{
		Name:     p.Name,
		Cases:    cases,
		CharData: helper.ReplaceNewLineAndTabToSpace(p.CharData),
	}
}

type caseComponent struct {
	Name     string `xml:"name,attr"`
	CharData string `xml:",chardata"`
}

func (c *caseComponent) toEntity() *entity.Case {
	return &entity.Case{
		Name:     c.Name,
		CharData: helper.ReplaceNewLineAndTabToSpace(c.CharData),
	}
}
