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
	case Sqlx:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  nil,
		}
	case Context:
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

func BuildErrWithOriginal(code errorCode, original error) Error {
	switch code.namespace {
	case Undefined:
		return &defaultError{
			errorCode: UndefinedErr,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	case FileIO:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	case SqlReader:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	case Query:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	case Sqlx:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	case Context:
		return &defaultError{
			errorCode: code,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	default:
		return &defaultError{
			errorCode: UndefinedErr,
			message:   nil,
			args:      map[string]string{},
			original:  original,
		}
	}
}
