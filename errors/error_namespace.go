package errors

type namespace int

const (
	Undefined namespace = iota
	SqlReader
	FileIO
)