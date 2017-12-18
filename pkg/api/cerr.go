package api

import (
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	ErrScope = "api"
)

var ErrClientUnknown = errx.New(ErrScope, "client-unknown")
var ErrUserUnknown = errx.New(ErrScope, "user-unknown")
var ErrDeviceUnknown = errx.New(ErrScope, "device-unknown")

var ErrInvalidSignature = errx.New(ErrScope, "invalid-signature")

var ErrSessionNotFound = errx.New(ErrScope, "session-not-found")

var ErrUserAlreadyExist = errx.New(ErrScope, "user-already-exist")

var ErrInvalidRequest = errx.New(ErrScope, "invalid-request")
var ErrTokenExpired = errx.New(ErrScope, "token-expired")
var ErrTokenInvalid = errx.New(ErrScope, "token-invalid")
var ErrAccessForbiden = errx.New(ErrScope, "access-forbiden")

var ErrAlreadyConfirmed = errx.New(ErrScope, "already-confirmed")

var ErrDeleteCurrentDevice = errx.New(ErrScope, "delete-current-device")
