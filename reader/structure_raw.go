package reader

import (
	"encoding/xml"
	"sql-mapper/entity"
	"sql-mapper/errors"
	"sql-mapper/helper"
)

type dmlBodyRaw struct {
	XMLName    xml.Name     `xml:"Body"`
	SelectRaws []*selectRaw `xml:"Select"`
	InputRaws  []*insertRaw `xml:"Insert"`
	UpdateRaws []*updateRaw `xml:"Update"`
	DeleteRaws []*deleteRaw `xml:"Delete"`
}

func (b dmlBodyRaw) toEntity(absFilePath string) (*entity.DMLBody, errors.Error) {
	var s []*entity.Select
	for _, sql := range b.SelectRaws {
		elem, err := sql.toEntity()
		if err != nil {
			return nil, err
		}
		s = append(s, elem)
	}

	var i []*entity.Insert
	for _, sql := range b.InputRaws {
		i = append(i, sql.toEntity())
	}

	var u []*entity.Update
	for _, sql := range b.UpdateRaws {
		u = append(u, sql.toEntity())
	}

	var d []*entity.Delete
	for _, sql := range b.DeleteRaws {
		d = append(d, sql.toEntity())
	}

	return &entity.DMLBody{
		AbsFilePath: absFilePath,
		Selects:     s,
		Inserts:     i,
		Updates:     u,
		Deletes:     d,
	}, nil
}

type insertRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s insertRaw) toEntity() *entity.Insert {

	return &entity.Insert{
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
}

type updateRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s updateRaw) toEntity() *entity.Update {

	return &entity.Update{
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
}

type deleteRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s deleteRaw) toEntity() *entity.Delete {

	return &entity.Delete{
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
}

type selectRaw struct {
	XMLName  xml.Name   `xml:"Select"`
	CharData string     `xml:",chardata"`
	Name     string     `xml:"name,attr"`
	PartRaws []*partRaw `xml:"Part"`
}

func (s *selectRaw) toEntity() (*entity.Select, errors.Error) {
	if helper.IsBlank(s.CharData) && len(s.PartRaws) == 0 {
		return nil, errors.BuildBasicErr(errors.ParseQueryErr)
	}

	part := []*entity.Part{}
	if helper.IsBlank(s.CharData) {
		for _, raw := range s.PartRaws {
			part = append(part, raw.toEntity())
		}
	}

	return &entity.Select{
		RawSql:    s.CharData,
		SimpleSql: !helper.IsBlank(s.CharData),
		Parts:     part,
		Name:      s.Name,
	}, nil
}

type partRaw struct {
	Name     string    `xml:"name,attr"`
	CharData string    `xml:",chardata"`
	CaseRaws []caseRaw `xml:"Case"`
}

func (p *partRaw) toEntity() *entity.Part {

	cases := []*entity.Case{}
	for _, raw := range p.CaseRaws {
		cases = append(cases, raw.toEntity())
	}

	return &entity.Part{
		Name:     p.Name,
		Cases:    cases,
		CharData: helper.ReplaceNewLineAndTabToSpace(p.CharData),
	}
}

type caseRaw struct {
	Name     string `xml:"name,attr"`
	CharData string `xml:",chardata"`
}

func (c *caseRaw) toEntity() *entity.Case {
	return &entity.Case{
		Name:     c.Name,
		CharData: helper.ReplaceNewLineAndTabToSpace(c.CharData),
	}
}
