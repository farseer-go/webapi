package check

import (
	"github.com/farseer-go/fs/exception"
	"strings"
)

// IsEmpty 当val为empty时，抛出异常
func IsEmpty(val string, statusCode int, err string) {
	if strings.TrimSpace(val) == "" {
		exception.ThrowWebException(statusCode, err)
	}
}

// IsTrue 当val为true时，抛出异常
func IsTrue(val bool, statusCode int, err string) {
	if val {
		exception.ThrowWebException(statusCode, err)
	}
}

// IsFalse 当val为false时，抛出异常
func IsFalse(val bool, statusCode int, err string) {
	if !val {
		exception.ThrowWebException(statusCode, err)
	}
}
