package errx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErr(t *testing.T) {
	errs := []Err{New("test", "test-reason"), NewWithFmt("test", "test-reason", "message-%d", 10), NewWithErr("test", errors.New("error"))}
	msgs := []string{"test.test-reason", "test.test-reason: message-10", "test.system: error"}

	json := []map[string]string{
		{"scope": "test", "reason": "test-reason"},
		{"scope": "test", "reason": "test-reason", "message": "message-10"},
		{"scope": "test", "reason": "system", "message": "error"},
	}

	for i, e := range errs {
		assert.Equal(t, msgs[i], e.Error())
		assert.Equal(t, json[i], e.Json())
	}
}

func TestErrJson(t *testing.T) {
	errs := []error{New("test", "test-reason"), NewWithErr("test", errors.New("error")), errors.New("error2")}
	json := []map[string]string{
		{"scope": "test", "reason": "test-reason"},
		{"scope": "test", "reason": "system", "message": "error"},
		{"scope": "test", "reason": "system", "message": "error2"},
	}

	for i, e := range errs {
		assert.Equal(t, json[i], Error2Json(e, "test"))
	}

}
