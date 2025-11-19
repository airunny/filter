package operation

import (
	"context"
	"runtime/debug"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variable"
	"github.com/stretchr/testify/assert"
)

type mockOperation struct {
	name string
}

func (m mockOperation) Name() string {
	return m.name
}

func (m mockOperation) PrepareValue(value interface{}) (interface{}, error) {
	return value, nil
}

func (m mockOperation) Run(ctx context.Context, variable variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	return true, nil
}

func TestRegister(t *testing.T) {
	f := func() {
		Register(nil)
	}
	funcDidPanic, panicValue, _ := didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register a nil Operation" {
		t.Fatalf("panic error got %s want cannot register a nil Operation", panicValue)
	}

	f = func() {
		Register(mockOperation{})
	}
	funcDidPanic, panicValue, _ = didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register Operation with empty string result for Name()" {
		t.Fatalf("panic error got %s want cannot register Operation with empty string result for Name()", panicValue)
	}

	op := mockOperation{
		name: "gt",
	}
	Register(op)
	got, ok := Get("gt")
	assert.True(t, ok)
	if got != op {
		t.Fatalf("Register(%v) want %v got %v", op, op, got)
	}

	got, ok = Get("lt")
	assert.False(t, ok)
	assert.Nil(t, got)
}

type PanicTestFunc func()

func didPanic(f PanicTestFunc) (bool, any, string) {
	didPanic := false
	var message any
	var stack string
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
				stack = string(debug.Stack())
			}
		}()

		// call the target function
		f()
	}()

	return didPanic, message, stack
}
