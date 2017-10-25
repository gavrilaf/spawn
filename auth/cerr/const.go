package cerr

import (
	"github.com/gavrilaf/go-auth/errors"
)

const (
	Scope         = "auth"
	ReasonDefault = "default"
)

var ClientUnknown = errors.NewErr(Scope, "client-unknown")
var UserUnknown = errors.NewErr(Scope, "user-unknown")
var InvalidSignature = errors.NewErr(Scope, "invalid-signature")

var SessionNotFound = errors.NewErr(Scope, "session-not-found")

var UserAlreadyExist = errors.NewErr(Scope, "user-already-exist")

var InvalidRequest = errors.NewErr(Scope, "invalid-request")
var TokenExpired = errors.NewErr(Scope, "token-expired")
var TokenInvalid = errors.NewErr(Scope, "token-invalid")
var AccessForbiden = errors.NewErr(Scope, "access-forbiden")
