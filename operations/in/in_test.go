package in

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/variables"
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

func TestIn(t *testing.T) {
	err := errors.New("value not found")
	cc := cache.NewCache()
	cases := []struct {
		Variable variables.Variable
		Value    interface{}
		Target   interface{}
		Data     interface{}
		Result   bool
		Err      error
		ParseErr error
	}{
		// int
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
				err:   err,
			},
			Value:  "[1,2]",
			Target: []interface{}{float64(1), float64(2)},
			Err:    err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:  "[1,2]",
			Target: []interface{}{float64(1), float64(2)},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "1",
			},
			Value:  `["1","2"]`,
			Target: []interface{}{"1", "2"},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "2",
			},
			Value:  `["1","2"]`,
			Target: []interface{}{"1", "2"},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "3",
			},
			Value:  `["1","2"]`,
			Target: []interface{}{"1", "2"},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "1",
			},
			Value:  `1,2`,
			Target: []interface{}{"1", "2"},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "2",
			},
			Value:  `1,2`,
			Target: []interface{}{"1", "2"},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "3",
			},
			Value:  `1,2`,
			Target: []interface{}{"1", "2"},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "3",
			},
			Value:  `["1","2"]`,
			Target: []interface{}{"1", "2"},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "2",
			},
			Value:    `[]`,
			ParseErr: emptyElementErr,
			Result:   false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"1", "2"},
			},
			Value:  `["1","2"]`,
			Target: []interface{}{"1", "2"},
			Result: true,
		},
	}

	op, ok := operations.Get("in")
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, "in", op.Name())

	for index, tt := range cases {
		value, err := op.PrepareValue(tt.Value)
		if tt.ParseErr != nil {
			assert.Equal(t, tt.ParseErr, err)
		} else {
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(value, tt.Target), index)
		}

		result, err := op.Run(context.Background(), tt.Variable, value, tt.Data, cc)
		if tt.Err != nil {
			assert.Equal(t, tt.Err, err)
			assert.Equal(t, tt.Result, result)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.Result, result, index)
		}
	}
}
