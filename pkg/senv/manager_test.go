package senv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevelopmentEnvironment(t *testing.T) {
	env := GetEnvironment()

	assert.Equal(t, env.GetName(), "Development")

	// TODO: Add missign checks & tests
}
