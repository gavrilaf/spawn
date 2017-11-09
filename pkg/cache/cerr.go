package cache

import (
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	errScope       = "auth"
	reasonNotFound = "key-not-found"
)

func errNotFound(key string) error {
	return errx.NewWithFmt(errScope, reasonNotFound, "Key %v not found", key)
}
