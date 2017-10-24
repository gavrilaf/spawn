package errors

import (
	"fmt"
)

type Err struct {
	scope  string
	reason string
	msg    string
}

func NewErr(scope string, reason string) Err {
	return Err{scope: scope, reason: reason, msg: ""}
}

func NewErrFormat(scope string, reason string, format string, args ...interface{}) Err {
	return Err{scope: scope, reason: reason, msg: fmt.Sprintf(format, args...)}
}

func (e Err) Error() string {
	return fmt.Sprintf("%v.%v: %v", e.scope, e.reason, e.msg)
}

func (e Err) Scope() string {
	return e.scope
}

func (e Err) Reason() string {
	return e.reason
}

func (e Err) Message() string {
	return e.msg
}
