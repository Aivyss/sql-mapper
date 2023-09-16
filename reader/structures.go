package reader

type Body struct {
	Selects []Select
	Inputs  []Input
	Updates []Update
	Deletes []Delete
	Creates []Create
	Drops   []Drop
}

type CommonFields struct {
	Sql      string
	Name     string
	FilePath string
}

type Select struct {
	List bool
	CommonFields
}

type Input struct {
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
