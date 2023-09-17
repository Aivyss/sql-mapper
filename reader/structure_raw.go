package reader

import (
	"encoding/xml"
	"sql-mapper/entity"
	"strconv"
)

type dmlBodyRaw struct {
	XMLName    xml.Name    `xml:"Body"`
	SelectRaws []selectRaw `xml:"Select"`
	InputRaws  []insertRaw `xml:"Insert"`
	UpdateRaws []updateRaw `xml:"Update"`
	DeleteRaws []deleteRaw `xml:"Delete"`
}

func (b dmlBodyRaw) toEntity(absFilePath string) *entity.DMLBody {
	var s []entity.Select
	for _, sql := range b.SelectRaws {
		s = append(s, *sql.toEntity())
	}

	var i []entity.Insert
	for _, sql := range b.InputRaws {
		i = append(i, *sql.toEntity())
	}

	var u []entity.Update
	for _, sql := range b.UpdateRaws {
		u = append(u, *sql.toEntity())
	}

	var d []entity.Delete
	for _, sql := range b.DeleteRaws {
		d = append(d, *sql.toEntity())
	}

	return &entity.DMLBody{
		AbsFilePath: absFilePath,
		Selects:     s,
		Inserts:     i,
		Updates:     u,
		Deletes:     d,
	}
}

type selectRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
	List string `xml:"list,attr"` // "true" or "false"
}

func (s selectRaw) toEntity() *entity.Select {
	listBool, err := strconv.ParseBool(s.List)
	if err != nil {
		listBool = false
	}

	return &entity.Select{
		List: listBool,
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
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

type createRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s createRaw) toEntity() *entity.Create {

	return &entity.Create{
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
}

type dropRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s dropRaw) toEntity() *entity.Drop {

	return &entity.Drop{
		CommonFields: entity.CommonFields{
			Sql:  s.Sql,
			Name: s.Name,
		},
	}
}
