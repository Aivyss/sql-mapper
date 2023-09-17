package entity

type DdlBody struct {
	Creates []Create
	Drops   []Drop
}

type Create struct {
	CommonFields
}

type Drop struct {
	CommonFields
}
