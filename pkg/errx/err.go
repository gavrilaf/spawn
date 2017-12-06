package errx

import (
	"fmt"
)

const (
	ScopeUnknown = "unknown"

	ReasonSystem         = "system"
	ReasonNotFound       = "not-found"
	ReasonNotImplemented = "not-implemented"
	ReasonInvalidFormat  = "invalid-format"
	ReasonEnvironment    = "environment"
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

func New(scope string, reason string) Err {
	return Err{scope: scope, reason: reason, msg: ""}
}

func NewFmt(scope string, reason string, format string, args ...interface{}) Err {
	return Err{scope: scope, reason: reason, msg: fmt.Sprintf(format, args...)}
}

func WrapErr(scope string, err error) Err {
	return Err{scope: scope, reason: ReasonSystem, msg: err.Error()}
}

func ErrNotFound(scope string, format string, args ...interface{}) Err {
	return NewFmt(scope, ReasonNotFound, format, args...)
}

func ErrEnvironment(scope string, format string, args ...interface{}) Err {
	return NewFmt(scope, ReasonEnvironment, format, args...)
}

func ErrKeyNotFound(scope string, key string) Err {
	return NewFmt(scope, ReasonNotFound, "Key %v not found", key)
}

////////////////////////////////////////////////////////////////////////////////////////////

// GetErrorReason returns Scope & Reason for error. For system error return (ScopeUnknown, ReasonSystem)
func GetErrorReason(err error) (string, string) {
	switch err2 := err.(type) {
	case Err:
		return err2.Scope(), err2.Reason()
	default:
		return ScopeUnknown, ReasonSystem
	}
}

////////////////////////////////////////////////////////////////////////////////////////////

// Encodable
func (e Err) ToMap() map[string]interface{} {
	j := map[string]interface{}{"scope": e.scope, "reason": e.reason}
	if e.msg != "" {
		j["message"] = e.msg
	}
	return j
}

func Error2Map(e error, defScope string) map[string]interface{} {
	switch e2 := e.(type) {
	case Err:
		return e2.ToMap()
	default:
		return WrapErr(defScope, e).ToMap()
	}
}
