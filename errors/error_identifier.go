package errors

type errorIdentifier int

const (
	UndefinedID errorIdentifier = iota
	ReadBodyErrID
	FileReadErrID
	FindQueryErrID
	ExecuteQueryErrID
	BeginTxErrID
	CommitTxErrID
	NoTxErrID
	ParseSqlErrID
	NotFoundQueryClientErrID
	RegisterQueryClientErrID
	BootstrapErrID
	WrongReadOnlySettingErrID
	StartTransactionErrID
	UnknownTransactionErrID
	TransactionClientSideErrID
	TransactionClientSidePanicErrID
)
