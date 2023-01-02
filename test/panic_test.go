package test

import (
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPanic(t *testing.T) {
	assert.Panics(t, func() {
		webapi.RegisterPOST("/", func() {})
	})
}
