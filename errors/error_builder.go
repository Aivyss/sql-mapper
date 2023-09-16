package errors

func BuildBasicErr(code errorCode) Error {
	switch code.namespace {
	case Undefined:
		return &defaultError{
			errorCode: UndefinedErr,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	case FileIO:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	case SqlReader:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	case Query:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	default:
		return &defaultError{
			errorCode: UndefinedErr,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	}
}
