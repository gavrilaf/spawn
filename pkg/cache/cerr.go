package cache

import (
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	errScope             = "auth"
	reasonNotFound       = "key-not-found"
	reasonNotImplemented = "not-implemented"
)

func errNotImplemented(s string) error {
	return errx.NewWithFmt(errScope, reasonNotFound, "Functoion %v not implemented", s)
}

func errNotFound(key string) error {
	return errx.NewWithFmt(errScope, reasonNotFound, "Key %v not found", key)
}
