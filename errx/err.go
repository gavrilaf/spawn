package errx

import (
	"fmt"
)

const (
	ReasonDefault = "default"
	ReasonSystem  = "system"
)

////////////////////////////////////////////////////////////////////////////////////////////

type Err struct {
	scope  string
	reason string
	msg    string
}

////////////////////////////////////////////////////////////////////////////////////////////

func (e Err) Scope() string {
	return e.scope
}

func (e Err) Reason() string {
	return e.reason
}

func (e Err) Message() string {
	return e.msg
}

////////////////////////////////////////////////////////////////////////////////////////////

func (e Err) Error() string {
	if e.msg == "" {
		return fmt.Sprintf("%v.%v", e.scope, e.reason)
	}
	return fmt.Sprintf("%v.%v: %v", e.scope, e.reason, e.msg)
}

////////////////////////////////////////////////////////////////////////////////////////////

func (e Err) Json() map[string]string {
	j := map[string]string{"scope": e.scope, "reason": e.reason}
	if e.msg != "" {
		j["message"] = e.msg
	}
	return j
}

func Error2Json(e error, defScope string) map[string]string {
	switch e2 := e.(type) {
	case Err:
		return e2.Json()
	default:
		return NewWithErr(defScope, e).Json()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////

func New(scope string, reason string) Err {
	return Err{scope: scope, reason: reason, msg: ""}
}

func NewWithFmt(scope string, reason string, format string, args ...interface{}) Err {
	return Err{scope: scope, reason: reason, msg: fmt.Sprintf(format, args...)}
}

func NewWithErr(scope string, err error) Err {
	return Err{scope: scope, reason: ReasonSystem, msg: err.Error()}
}

////////////////////////////////////////////////////////////////////////////////////////////
