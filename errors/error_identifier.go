package errors

type errorIdentifier int

const (
	UndefinedID errorIdentifier = iota
	ReadBodyErrID
	FileReadErrID
)