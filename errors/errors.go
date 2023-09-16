package errors

import "fmt"

type Error interface {
	Code() *errorCode
	Msg() *string
	Args() *map[string]string
	error
}

type defaultError struct {
	errorCode
	message  *string
	args     map[string]string
	original error
}

func (e *defaultError) Code() *errorCode {
	return &e.errorCode
}

func (e *defaultError) Msg() *string {
	return e.message
}

func (e *defaultError) Args() *map[string]string {
	return &e.args
}

func (e *defaultError) Error() string {
	if e.original != nil {
		return e.original.Error()
	}

	var msg string
	if e.message == nil {
		msg = e.defaultMessage
	} else {
		msg = *e.message
	}
	return fmt.Sprintf(
		"[Msg] %v\n[code-identifier] %v\n[code-namespace] %v\n[args] %v",
		msg, e.identifier, e.namespace, e.args,
	)
}
