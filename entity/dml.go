package entity

type DMLBody struct {
	AbsFilePath string
	Selects     []*Select
	Inserts     []*Insert
	Deletes     []*Delete
	Updates     []*Update
}

type Case struct {
	CharData string
	Name     string
}

type Part struct {
	Name     string
	CharData string
	Cases    []*Case
}

type Select struct {
	Name      string
	RawSql    string
	Parts     []*Part
	SimpleSql bool
}

type CommonFields struct {
	Sql  string
	Name string
}

type Insert struct {
	CommonFields
}

type Update struct {
	CommonFields
}

type Delete struct {
	CommonFields
}
