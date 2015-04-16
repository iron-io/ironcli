package examples

import (
	"errors"
	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/arschles/assert"
	"testing"
)

func TestNil(t *testing.T) {
	assert.Nil(t, nil, "nil")
	assert.NotNil(t, "abc", "string")
}

func TestBooleans(t *testing.T) {
	assert.True(t, true, "boolean true")
	assert.False(t, false, "boolean false")
}

func TestEqual(t *testing.T) {
	s1 := struct {
		a string
		b int
	}{"testString", 1}
	s2 := struct {
		a string
		b int
	}{"testString", 1}
	assert.Equal(t, s1, s2, "anonymous struct")
}

func TestErrors(t *testing.T) {
	err1 := errors.New("this is an error")
	var err2 error = nil
	assert.Err(t, err1, errors.New("this is an error"))
	assert.NoErr(t, err2)
	assert.ExistsErr(t, err1, "valid error")
}
