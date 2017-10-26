package cerr

import (
	"github.com/gavrilaf/go-auth/errx"
)

const (
	Scope = "auth"
)

var ClientUnknown = errx.New(Scope, "client-unknown")
var UserUnknown = errx.New(Scope, "user-unknown")
var DeviceUnknown = errx.New(Scope, "device-unknown")

var InvalidSignature = errx.New(Scope, "invalid-signature")

var SessionNotFound = errx.New(Scope, "session-not-found")

var UserAlreadyExist = errx.New(Scope, "user-already-exist")

var InvalidRequest = errx.New(Scope, "invalid-request")
var TokenExpired = errx.New(Scope, "token-expired")
var TokenInvalid = errx.New(Scope, "token-invalid")
var AccessForbiden = errx.New(Scope, "access-forbiden")
