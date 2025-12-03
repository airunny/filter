package assignment

import (
	"context"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAssignment struct {
	name string
	err  error
}

func (m mockAssignment) Name() string {
	return m.name
}

func (m mockAssignment) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return value, nil
}

func (m mockAssignment) Run(ctx context.Context, data interface{}, key string, val interface{}) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestRegister(t *testing.T) {
	f := func() {
		Register(nil)
	}
	funcDidPanic, panicValue, _ := didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register a nil Assignment" {
		t.Fatalf("panic error got %s want cannot register a nil Assignment", panicValue)
	}

	f = func() {
		Register(mockAssignment{})
	}
	funcDidPanic, panicValue, _ = didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register Assignment with empty string result for Name()" {
		t.Fatalf("panic error got %s want cannot register Assignment with empty string result for Name()", panicValue)
	}

	op := mockAssignment{
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
