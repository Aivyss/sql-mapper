package errors

type errorIdentifier int

const (
	UndefinedID errorIdentifier = iota
	ReadBodyErrID
	FileReadErrID
	FindQueryErrID
	DuplicatedIdentifierErrID
	FindQueryMapErrID
	ExecuteQueryErrID
	BeginTxErrID
	CommitTxErrID
	NoTxErrID
	ParseSqlErrID
	NotFoundQueryClientErrID
	RegisterQueryClientErrID
)
