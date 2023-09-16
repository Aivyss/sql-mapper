package errors

var (
	UndefinedErr     = newErrorCode(UndefinedID, Undefined, "undefined error")
	ReadBodyErr      = newErrorCode(ReadBodyErrID, SqlReader, "fail to unmarshal xml file")
	FileReadErr      = newErrorCode(FileReadErrID, FileIO, "fail to file read")
	FindQueryErr     = newErrorCode(FindQueryErrID, Query, "fail to find the query")
	RegisterQueryErr = newErrorCode(RegisterQueryErrID, Query, "fail to register the query set")
)

type errorCode struct {
	identifier errorIdentifier
	namespace
	defaultMessage string
}

func newErrorCode(identifier errorIdentifier, namespace namespace, message string) errorCode {
	return errorCode{
		identifier:     identifier,
		namespace:      namespace,
		defaultMessage: message,
	}
}
