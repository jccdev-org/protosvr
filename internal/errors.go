package internal

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func WrapError(err error) error {
	_, filename, line, _ := runtime.Caller(1)
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		unwrapped = err
	}
	existingStack := WrappedErrorStack(err)
	return fmt.Errorf("%w[at] %s:%d %s", unwrapped, filename, line, existingStack)
}

// WrapErrorN n will pass to runtime.Caller skip parameter
// use n to customize how many stack frames runtime.Caller will skip
func WrapErrorN(err error, n int) error {
	_, filename, line, _ := runtime.Caller(n)
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		unwrapped = err
	}
	existingStack := WrappedErrorStack(err)
	return fmt.Errorf("%w[at] %s:%d %s", unwrapped, filename, line, existingStack)
}

func WrappedErrorMsg(err error) string {
	errMsg := err.Error()
	matchIx := strings.Index(errMsg, "[at]")
	if matchIx != -1 {
		errMsg = errMsg[:matchIx]
	}
	return errMsg
}

func WrappedErrorStack(err error) string {
	errMsg := err.Error()
	existingStack := ""
	matchIx := strings.Index(errMsg, "[at]")
	if matchIx != -1 {
		existingStack = errMsg[matchIx:]
	}
	return existingStack
}

func PrettyPrintError(err error) string {
	return fmt.Sprintf("[Error] %s", strings.Replace(err.Error(), "[at]", "\n\t[at]", -1))
}
