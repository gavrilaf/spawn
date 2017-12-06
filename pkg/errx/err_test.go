package errx

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErr_Error(t *testing.T) {
	errs := []Err{
		New("test", "test-reason"),
		NewFmt("test", "test-reason", "message-%d", 10),
		WrapErr("test", errors.New("error")),
		ErrNotFound("test", "err: %d", 1),
		ErrEnvironment("test2", "err"),
	}

	msgs := []string{
		"test.test-reason",
		"test.test-reason: message-10",
		"test.system: error",
		"test.not-found: err: 1",
		"test2.environment: err",
	}

	json := []map[string]interface{}{
		{"scope": "test", "reason": "test-reason"},
		{"scope": "test", "reason": "test-reason", "message": "message-10"},
		{"scope": "test", "reason": "system", "message": "error"},
		{"scope": "test", "reason": "not-found", "message": "err: 1"},
		{"scope": "test2", "reason": "environment", "message": "err"},
	}

	for i, e := range errs {
		assert.Equal(t, msgs[i], e.Error())
		assert.Equal(t, json[i], e.ToMap())
	}
}

func TestErr_ToMap(t *testing.T) {
	errs := []error{
		New("test", "test-reason"),
		WrapErr("test", errors.New("error")),
		errors.New("error2"),
	}

	json := []map[string]interface{}{
		{"scope": "test", "reason": "test-reason"},
		{"scope": "test", "reason": "system", "message": "error"},
		{"scope": "test", "reason": "system", "message": "error2"},
	}

	for i, e := range errs {
		assert.Equal(t, json[i], Error2Map(e, "test"))
	}
}

func Test_GetErrorReason(t *testing.T) {
	errs := []error{
		New("test1", "test-reason"),
		ErrNotFound("test2", "err: %d", 1),
		ErrEnvironment("test3", "err"),
		fmt.Errorf("simple"),
	}

	expected := []struct {
		scope  string
		reason string
	}{
		{"test1", "test-reason"},
		{"test2", "not-found"},
		{"test3", "environment"},
		{"unknown", "system"},
	}

	for i, e := range errs {
		scope, reason := GetErrorReason(e)
		assert.Equal(t, expected[i].scope, scope)
		assert.Equal(t, expected[i].reason, reason)
	}
}
