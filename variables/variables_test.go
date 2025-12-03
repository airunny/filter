package variables

import (
	"context"
	"runtime/debug"
	"testing"

	"github.com/airunny/filter/cache"
	"github.com/stretchr/testify/assert"
)

type mockVariable struct {
	name  string
	value interface{}
}

func (m mockVariable) Name() string {
	return m.name
}

func (m mockVariable) Cacheable() bool {
	return true
}

func (m mockVariable) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	return m.value, nil
}

func TestRegister(t *testing.T) {
	f := func() {
		Register(nil)
	}
	funcDidPanic, panicValue, _ := didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register a nil variable builder" {
		t.Fatalf("panic error got %s want cannot register a nil variable builder", panicValue)
	}

	f = func() {
		Register(NewSimpleVariable(&mockVariable{}))
	}
	funcDidPanic, panicValue, _ = didPanic(f)
	if !funcDidPanic {
		t.Fatalf("func should panic\n\tPanic value:\t%#v", panicValue)
	}
	if panicValue != "cannot register variable builder with empty string result for Name()" {
		t.Fatalf("panic error got %s want cannot register variable builder with empty string result for Name()", panicValue)
	}

	op := &mockVariable{
		name: "mock",
	}
	Register(NewSimpleVariable(op))
	got, ok := Get("mock")
	assert.True(t, ok)
	if got != op {
		t.Fatalf("Register(%v) want %v got %v", op, op, got)
	}
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
