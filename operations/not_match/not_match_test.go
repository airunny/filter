package not_match

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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

func TestNotMatch(t *testing.T) {
	cc := cache.NewCache()
	err := errors.New("value not found")
	cases := []struct {
		Variable    variables.Variable
		Value       interface{}
		ParsedValue func() interface{}
		ParsedErr   error
		Data        interface{}
		Result      bool
		ResultErr   error
	}{
		// err
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value: "",
			ParsedValue: func() interface{} {
				return ""
			},
			Result:    false,
			ResultErr: err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value:       1,
			ParsedValue: nil,
			ParsedErr:   ErrInvalidOperationValue,
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:       1,
			ParsedValue: nil,
			ParsedErr:   ErrInvalidOperationValue,
			Result:      false,
			ResultErr:   ErrInvalidVariableValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value:       1,
			ParsedValue: nil,
			ParsedErr:   ErrInvalidOperationValue,
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value:       "//",
			ParsedValue: nil,
			ParsedErr:   fmt.Errorf("[%s] operation value is not a valid regexp [%s]", Name, "//"),
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value:       "/(aa/",
			ParsedValue: nil,
			ParsedErr:   fmt.Errorf("[%s] operation value is not a valid regexp [%s]", Name, "/(aa/"),
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		// string
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: "lang",
			ParsedValue: func() interface{} {
				return "lang"
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: "go",
			ParsedValue: func() interface{} {
				return "go"
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: "Lang",
			ParsedValue: func() interface{} {
				return "Lang"
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: "",
			ParsedValue: func() interface{} {
				return ""
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value: "",
			ParsedValue: func() interface{} {
				return ""
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value: "go",
			ParsedValue: func() interface{} {
				return "go"
			},
			Result: true,
		},
		// regex
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "test@example.com",
			},
			Value: "/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/",
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return reg
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "user.name+tag@domain.co.uk",
			},
			Value: "/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/",
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return reg
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "invalid.email",
			},
			Value: "/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/",
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return reg
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "missing@tld.",
			},
			Value: "/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/",
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return reg
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "@domain.com",
			},
			Value: "/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/",
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return reg
			},
			Result: true,
		},
	}

	op, ok := operations.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, op)
	assert.Equal(t, Name, op.Name())

	for index, tt := range cases {
		value, err := op.PrepareValue(tt.Value)
		if tt.ParsedErr != nil {
			assert.Equal(t, tt.ParsedErr, err)
			assert.Nil(t, value)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.ParsedValue(), value, index)
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
