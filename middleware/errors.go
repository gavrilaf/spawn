package middleware

import (
	"github.com/gavrilaf/go-auth/errors"
)

const (
	ErrScope         = "auth"
	ErrReasonDefault = "default"
)

var errInvalidRequest = errors.NewErr(ErrScope, "invalid-request")
var errTokenExpired = errors.NewErr(ErrScope, "token-expired")
var errTokenInvalid = errors.NewErr(ErrScope, "token-invalid")
var errAccessForbiden = errors.NewErr(ErrScope, "access-forbiden")
