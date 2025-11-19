package greater_than_equal

import (
	"context"
	"errors"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/variable"
	"github.com/stretchr/testify/assert"
)

type mockVariable struct {
	name      string
	cacheable bool
	err       error
	value     interface{}
}

func (s mockVariable) Name() string {
	return s.name
}

func (s mockVariable) Cacheable() bool {
	return s.cacheable
}

func (s mockVariable) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.value, nil
}

func TestEqual(t *testing.T) {
	err := errors.New("value not found")
	cc := cache.NewCache()
	cases := []struct {
		Variable variable.Variable
		Value    interface{}
		Data     interface{}
		Result   bool
		Err      error
	}{
		// int
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:  0,
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:  1,
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:  2,
			Result: false,
		},
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Err: err,
		},
		// string
		{
			Variable: mockVariable{
				name:  "mock",
				value: "1",
			},
			Value:  "0",
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "0",
			},
			Value:  "1",
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "1",
			},
			Value:  "1",
			Result: true,
		},
		// float
		{
			Variable: mockVariable{
				name:  "mock",
				value: 10.19,
			},
			Value:  10.18,
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 10.18,
			},
			Value:  10.18,
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 10.18,
			},
			Value:  10.19,
			Result: false,
		},
		// string list
		{
			Variable: mockVariable{
				name:  "mock",
				value: []string{"1", "2"},
			},
			Value:  []string{"1", "2"},
			Result: true,
		},
		// int list
		{
			Variable: mockVariable{
				name:  "mock",
				value: []int{1, 2},
			},
			Value:  []int{1, 2},
			Result: true,
		},
		// interface list
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{1, 2},
			},
			Value:  []interface{}{1, 2},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{1, "golang"},
			},
			Value:  []interface{}{1, "golang"},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: nil,
			},
			Value:  nil,
			Result: true,
		},
	}

	op, ok := operation.Get(">=")
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, ">=", op.Name())

	op, ok = operation.Get("gte")
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, "gte", op.Name())

	for index, tt := range cases {
		value, err := op.PrepareValue(tt.Value)
		assert.Nil(t, err)
		assert.Equal(t, value, tt.Value)

		result, err := op.Run(context.Background(), tt.Variable, value, tt.Data, cc)
		if tt.Err != nil {
			assert.Equal(t, tt.Err, err, index)
			assert.Equal(t, tt.Result, result, index)
		} else {
			assert.Nil(t, err, index)
			assert.Equal(t, tt.Result, result, index)
		}
	}
}
