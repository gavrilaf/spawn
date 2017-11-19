package profile

import (
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	errScope = "auth"
)

var errAlreadyConfirmed = errx.New(errScope, "device-already-confirmed")
var errInvalidConfirm = errx.New(errScope, "invalid-confirm-code")
