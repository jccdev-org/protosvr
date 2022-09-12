package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWrapError(t *testing.T) {
	// arrange and act
	orig := errors.New("original error")
	first := WrapError(orig)
	second := WrapError(first)
	third := WrapError(second)

	// assert
	assert.True(t, strings.HasPrefix(third.Error(), "original error"), "should starts with original error message")
	assert.Regexp(t, "\\[at\\].*\\/errors_test.go:13", third, "should contain line 13 in stack")
	assert.Regexp(t, "\\[at\\].*\\/errors_test.go:14", third, "should contain line 14 in stack")
	assert.Regexp(t, "\\[at\\].*\\/errors_test.go:15", third, "should contain line 15 in stack")
}

func TestWrappedErrorMsg(t *testing.T) {
	// arrange
	orig := errors.New("original error")
	first := WrapError(orig)
	second := WrapError(first)
	third := WrapError(second)

	// act
	msg := WrappedErrorMsg(third)

	// assert
	assert.Equal(t, "original error", msg)
}

func TestPrettyPrintError(t *testing.T) {
	// arrange
	orig := errors.New("original error")
	first := WrapError(orig)
	second := WrapError(first)
	third := WrapError(second)

	// act
	msg := PrettyPrintError(third)
	lines := strings.Split(msg, "\n")

	// assert
	assert.Regexp(t, "\\[Error\\] original error", lines[0])
	assert.Regexp(t, "\\t\\[at\\].*\\/errors_test.go:43", lines[1], "should contain line 43 in stack")
	assert.Regexp(t, "\\t\\[at\\].*\\/errors_test.go:42", lines[2], "should contain line 42 in stack")
	assert.Regexp(t, "\\t\\[at\\].*\\/errors_test.go:41", lines[3], "should contain line 41 in stack")
}
