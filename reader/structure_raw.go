package reader

import (
	"encoding/xml"
	"strconv"
)

type bodyRaw struct {
	XMLName    xml.Name    `xml:"Body"`
	SelectRaws []selectRaw `xml:"Select"`
	InputRaws  []inputRaw  `xml:"Input"`
	UpdateRaws []updateRaw `xml:"Update"`
	DeleteRaws []deleteRaw `xml:"Delete"`
	CreateRaws []createRaw `xml:"Create"`
	DropRaws   []dropRaw   `xml:"Drop"`
}

func (b bodyRaw) toEntity(filePath string) *Body {
	var s []Select
	for _, sql := range b.SelectRaws {
		s = append(s, *sql.toEntity(filePath))
	}

	var i []Input
	for _, sql := range b.InputRaws {
		i = append(i, *sql.toEntity(filePath))
	}

	var u []Update
	for _, sql := range b.UpdateRaws {
		u = append(u, *sql.toEntity(filePath))
	}

	var c []Create
	for _, sql := range b.CreateRaws {
		c = append(c, *sql.toEntity(filePath))
	}

	var d []Delete
	for _, sql := range b.DeleteRaws {
		d = append(d, *sql.toEntity(filePath))
	}

	var drop []Drop
	for _, sql := range b.DropRaws {
		drop = append(drop, *sql.toEntity(filePath))
	}

	return &Body{
		Selects: s,
		Inputs:  i,
		Updates: u,
		Creates: c,
		Deletes: d,
		Drops:   drop,
	}
}

type selectRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
	List string `xml:"list,attr"` // "true" or "false"
}

func (s selectRaw) toEntity(filePath string) *Select {
	listBool, err := strconv.ParseBool(s.List)
	if err != nil {
		listBool = false
	}

	return &Select{
		List: listBool,
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}

type inputRaw struct {
	Sql  string `xml:"Input"`
	Name string `xml:"name,attr"`
}

func (s inputRaw) toEntity(filePath string) *Input {

	return &Input{
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}

type updateRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s updateRaw) toEntity(filePath string) *Update {

	return &Update{
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}

type deleteRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s deleteRaw) toEntity(filePath string) *Delete {

	return &Delete{
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}

type createRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s createRaw) toEntity(filePath string) *Create {

	return &Create{
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}

type dropRaw struct {
	Sql  string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

func (s dropRaw) toEntity(filePath string) *Drop {

	return &Drop{
		CommonFields: CommonFields{
			Sql:      s.Sql,
			FilePath: filePath,
			Name:     s.Name,
		},
	}
}
