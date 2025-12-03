package any

import (
	"context"
	"errors"
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

func TestAny(t *testing.T) {
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
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"1", "2"},
				err:   err,
			},
			OperationValue: `["1","2"]`,
			ParsedValue:    []interface{}{"1", "2"},
			Result:         false,
			ResultErr:      err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"1", "2"},
			},
			OperationValue: `["1","2"]`,
			ParsedValue:    []interface{}{"1", "2"},
			Result:         true,
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
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{1},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{2},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{3},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"1"},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"2"},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"3"},
			},
			OperationValue: `[1,2]`,
			ParsedValue:    []interface{}{float64(1), float64(2)},
			Result:         false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"1.1"},
			},
			OperationValue: `[1.1,2.2]`,
			ParsedValue:    []interface{}{1.1, 2.2},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"2.2"},
			},
			OperationValue: `[1.1,2.2]`,
			ParsedValue:    []interface{}{1.1, 2.2},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"2.2"},
			},
			OperationValue: `["1.1","2.2"]`,
			ParsedValue:    []interface{}{"1.1", "2.2"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{"3.3"},
			},
			OperationValue: `["1.1","2.2"]`,
			ParsedValue:    []interface{}{"1.1", "2.2"},
			Result:         false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: []interface{}{[]interface{}{"1.1"}},
			},
			OperationValue: `[["1.1"],"2.2"]`,
			ParsedValue:    []interface{}{[]interface{}{"1.1"}, "2.2"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "2.2",
			},
			OperationValue: `[["1.1"],"2.2"]`,
			ParsedValue:    []interface{}{[]interface{}{"1.1"}, "2.2"},
			Result:         true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "3.3",
			},
			OperationValue: `[["1.1"],"2.2"]`,
			ParsedValue:    []interface{}{[]interface{}{"1.1"}, "2.2"},
			Result:         false,
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
			assert.Equal(t, tt.Result, result, index)
		}
	}
}
