package entity

type DMLBody struct {
	AbsFilePath string
	Selects     []Select
	Inserts     []Insert
	Updates     []Update
	Deletes     []Delete
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
