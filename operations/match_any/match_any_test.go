package match_any

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

func TestMatchAny(t *testing.T) {
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
				return []interface{}{""}
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
			ParsedErr:   ErrInvalidOperationElementValue,
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value:       "[]",
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
			Value:       `["//"]`,
			ParsedValue: nil,
			ParsedErr:   ErrEmptyOperationValue,
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "",
			},
			Value:       `["/(aa/"]`,
			ParsedValue: nil,
			ParsedErr:   fmt.Errorf("[%s] operation value invalid regexp [%s]", Name, "(aa"),
			Result:      false,
			ResultErr:   ErrInvalidOperationValue,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: 1,
			},
			Value: "golang",
			ParsedValue: func() interface{} {
				return []interface{}{"golang"}
			},
			ParsedErr: nil,
			Result:    false,
			ResultErr: ErrInvalidVariableValue,
		},
		// success
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: `lang`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang"}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: `["lang"]`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang"}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: `["lang","go"]`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang", "go"}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "golang",
			},
			Value: "lang,go",
			ParsedValue: func() interface{} {
				return []interface{}{"lang", "go"}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "java",
			},
			Value: `["lang","go"]`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang", "go"}
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "python",
			},
			Value: `["lang","go"]`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang", "go"}
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "go",
			},
			Value: `["lang","go"]`,
			ParsedValue: func() interface{} {
				return []interface{}{"lang", "go"}
			},
			Result: true,
		},
		// regex
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "test@example.com",
			},
			Value: `["/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{reg}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "user.name+tag@domain.co.uk",
			},
			Value: `["/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{reg}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "invalid.email",
			},
			Value: `["/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{reg}
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "missing@tld.",
			},
			Value: `["/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{reg}
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "@domain.com",
			},
			Value: `["/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{reg}
			},
			Result: false,
		},
		// string and regex
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "golang",
			},
			Value: `["golang","/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{"golang", reg}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "user.name+tag@domain.co.uk",
			},
			Value: `["golang","/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{"golang", reg}
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "@domain.com",
			},
			Value: `["golang","/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{"golang", reg}
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "邮箱验证",
				value: "@golang",
			},
			Value: `["golang","/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$/"]`,
			ParsedValue: func() interface{} {
				reg, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
				return []interface{}{"golang", reg}
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
