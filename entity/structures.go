package entity

type Body struct {
	AbsFilePath string
	Selects     []Select
	Inserts     []Insert
	Updates     []Update
	Deletes     []Delete
	Creates     []Create
	Drops       []Drop
}

type CommonFields struct {
	Sql  string
	Name string
}

type Select struct {
	List bool
	CommonFields
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

type Create struct {
	CommonFields
}

type Drop struct {
	CommonFields
}
