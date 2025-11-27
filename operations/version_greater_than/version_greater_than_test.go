package version_greater_than

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

func TestVersionGreaterThan(t *testing.T) {
	err := errors.New("value not found")
	cc := cache.NewCache()
	cases := []struct {
		Variable    variables.Variable
		Value       interface{}
		ParsedValue interface{}
		Data        interface{}
		Result      bool
		Err         error
	}{
		// int
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value:       0,
			ParsedValue: 0,
			Result:      false,
			Err:         err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v0",
			},
			Value:       "v0",
			ParsedValue: "v0",
			Result:      false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v0.0.1",
			},
			Value:       "v0.0.0",
			ParsedValue: "v0.0.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v0.0.2",
			},
			Value:       "v0.0.1",
			ParsedValue: "v0.0.1",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v0.0.1.0",
			},
			Value:       "v0.0.1",
			ParsedValue: "v0.0.1",
			Result:      false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v1.0.1",
			},
			Value:       "v1.0.0",
			ParsedValue: "v1.0.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v1.1.0",
			},
			Value:       "v1.0.0",
			ParsedValue: "v1.0.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v1.1.1",
			},
			Value:       "v1.0.0",
			ParsedValue: "v1.0.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v1.2.0",
			},
			Value:       "v1.0.0",
			ParsedValue: "v1.0.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v1.2.1",
			},
			Value:       "v1.2.0",
			ParsedValue: "v1.2.0",
			Result:      true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "v2.1.0",
			},
			Value:       "v1.0.0",
			ParsedValue: "v1.0.0",
			Result:      true,
		},
	}

	op, ok := operations.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, Name, op.Name())

	for _, tt := range cases {
		value, err := op.PrepareValue(tt.Value)
		assert.Nil(t, err)
		assert.Equal(t, value, tt.Value)

		result, err := op.Run(context.Background(), tt.Variable, value, tt.Data, cc)
		if tt.Err != nil {
			assert.Equal(t, tt.Err, err)
			assert.Equal(t, tt.Result, result)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.Result, result)
		}
	}
}
