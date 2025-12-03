package between

import (
	"context"
	"errors"
	"testing"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/variables"
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

func TestBetween(t *testing.T) {
	cc := cache.NewCache()
	err := errors.New("value not found")
	cases := []struct {
		Variable       variables.Variable
		OperationValue interface{}
		ParsedValue    interface{}
		ParsedErr      error
		Data           interface{}
		Result         bool
		ResultErr      error
	}{
		// err
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
				err:   err,
			},
			OperationValue: `[1,4]`,
			ParsedValue:    []interface{}{float64(1), float64(4)},
			Result:         false,
			ResultErr:      err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			OperationValue: `[]`,
			ParsedValue:    []interface{}{},
			ParsedErr:      ErrInvalidOperationValue,
			Result:         false,
			ResultErr:      ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			OperationValue: `[1]`,
			ParsedValue:    []interface{}{},
			ParsedErr:      ErrInvalidOperationValue,
			Result:         false,
			ResultErr:      ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			OperationValue: `[1,2,3]`,
			ParsedValue:    []interface{}{},
			ParsedErr:      ErrInvalidOperationValue,
			Result:         false,
			ResultErr:      ErrInvalidOperationValue,
		},
		// int
		{
			Variable: mockVariable{
				name:  "mock",
				value: 3,
			},
			OperationValue: `[1,4]`,
			ParsedValue:    []interface{}{float64(1), float64(4)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 0,
			},
			OperationValue: `[-1,4]`,
			ParsedValue:    []interface{}{float64(-1), float64(4)},
			Result:         true,
		},
		// float
		{
			Variable: mockVariable{
				name:  "mock",
				value: -1,
			},
			OperationValue: `[-2,0]`,
			ParsedValue:    []interface{}{float64(-2), float64(0)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2,
			},
			OperationValue: `[-3,-1]`,
			ParsedValue:    []interface{}{float64(-3), float64(-1)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2.0,
			},
			OperationValue: `[-3,-1]`,
			ParsedValue:    []interface{}{float64(-3), float64(-1)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2.0,
			},
			OperationValue: `[-3.0,-1.0]`,
			ParsedValue:    []interface{}{-3.0, -1.0},
			Result:         true,
		},
		// string
		{
			Variable: mockVariable{
				name:  "mock",
				value: 3,
			},
			OperationValue: `["1","4"]`,
			ParsedValue:    []interface{}{"1", "4"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 0,
			},
			OperationValue: `["-1","4"]`,
			ParsedValue:    []interface{}{"-1", "4"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -1,
			},
			OperationValue: `["-2","0"]`,
			ParsedValue:    []interface{}{"-2", "0"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2,
			},
			OperationValue: `["-3","-1"]`,
			ParsedValue:    []interface{}{"-3", "-1"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2.0,
			},
			OperationValue: `["-3","-1"]`,
			ParsedValue:    []interface{}{"-3", "-1"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: -2.0,
			},
			OperationValue: `["-3.0","-1.0"]`,
			ParsedValue:    []interface{}{"-3.0", "-1.0"},
			Result:         true,
		},
	}

	op, ok := operations.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, Name, op.Name())

	for index, tt := range cases {
		value, err := op.PrepareValue(tt.OperationValue)
		if tt.ParsedErr != nil {
			assert.Equal(t, tt.ParsedErr, err)
			assert.Nil(t, value)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.ParsedValue, value, index)
		}

		result, err := op.Run(context.Background(), tt.Variable, value, tt.Data, cc)
		if tt.ResultErr != nil {
			assert.Equal(t, tt.ResultErr, err)
			assert.Equal(t, tt.Result, result)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.Result, result)
		}
	}
}
