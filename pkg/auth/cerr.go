package auth

import (
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	errScope = "auth"
)

var errClientUnknown = errx.New(errScope, "client-unknown")
var errUserUnknown = errx.New(errScope, "user-unknown")
var errDeviceUnknown = errx.New(errScope, "device-unknown")

var errInvalidSignature = errx.New(errScope, "invalid-signature")

var errSessionNotFound = errx.New(errScope, "session-not-found")

var errUserAlreadyExist = errx.New(errScope, "user-already-exist")

var errInvalidRequest = errx.New(errScope, "invalid-request")
var errTokenExpired = errx.New(errScope, "token-expired")
var errTokenInvalid = errx.New(errScope, "token-invalid")
var errAccessForbiden = errx.New(errScope, "access-forbiden")
