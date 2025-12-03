package in_ip_range

import (
	"context"
	"errors"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
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

func TestInIpRange(t *testing.T) {
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
			},
			Value:     ``,
			ParsedErr: ErrEmptyOperationValueElement,
			ResultErr: ErrInvalidVariableValue,
		},
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value: `[]`,
			ParsedValue: func() interface{} {
				return nil
			},
			ParsedErr: ErrInvalidOperationValue,
			ResultErr: err,
		},
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value: `[1,2]`,
			ParsedValue: func() interface{} {
				return nil
			},
			ParsedErr: ErrInvalidOperationValueElement,
			ResultErr: err,
		},
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value: `[""]`,
			ParsedValue: func() interface{} {
				return nil
			},
			ParsedErr: ErrEmptyOperationValueElement,
			ResultErr: err,
		},
		{
			Variable: mockVariable{
				name: "mock",
				err:  err,
			},
			Value: `["192.0.2.0/24"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24")
				return vv
			},
			ResultErr: err,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.0.2.0",
			},
			Value: `["192.0.2.0/24"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.0.2.1",
			},
			Value: `["192.0.2.0/24"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.0.2.255",
			},
			Value: `["192.0.2.0/24"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.0.3.255",
			},
			Value: `["192.0.2.0/24"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24")
				return vv
			},
			Result: false,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.168.12.24",
			},
			Value: `["192.0.2.0/24","192.168.12.27/30"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24", "192.168.12.27/30")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.168.12.25",
			},
			Value: `["192.0.2.0/24","192.168.12.27/30"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24", "192.168.12.27/30")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.168.12.26",
			},
			Value: `["192.0.2.0/24","192.168.12.27/30"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24", "192.168.12.27/30")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.168.12.27",
			},
			Value: `["192.0.2.0/24","192.168.12.27/30"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24", "192.168.12.27/30")
				return vv
			},
			Result: true,
		},
		{
			Variable: mockVariable{
				name:  "mock",
				value: "192.168.12.28",
			},
			Value: `["192.0.2.0/24","192.168.12.27/30"]`,
			ParsedValue: func() interface{} {
				vv, _ := utils.IPRanges("192.0.2.0/24", "192.168.12.27/30")
				return vv
			},
			Result: false,
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
			assert.Nil(t, err, index)
			assert.Equal(t, tt.Result, result)
		}
	}
}
