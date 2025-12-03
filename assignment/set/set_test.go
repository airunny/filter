package set

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/liyanbing/filter/assignment"
	"github.com/stretchr/testify/assert"
)

type MockSetterData struct {
	Key   string
	Value interface{}
	err   error
}

func (d *MockSetterData) Set(key string, value interface{}) error {
	if d.err != nil {
		return d.err
	}

	d.Key = key
	d.Value = value
	return nil
}

type MockData struct {
	Key     string
	Value   interface{}
	Data    *MockData
	NilData interface{} `json:"nil_data"`
}

type SubData struct {
	Name string   `json:"name"`
	Data *SubData `json:"data"`
}

func TestSet(t *testing.T) {
	ass, ok := assignment.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, ass)
	assert.Equal(t, Name, ass.Name())

	var (
		err = errors.New("set value err")
		ctx = context.Background()
	)

	cases := []struct {
		Key          string
		Value        interface{}
		ParsedValue  interface{}
		Data         interface{}
		ResultErr    error
		ErrData      interface{}
		ExpectedData interface{}
	}{
		// nil data
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data:        nil,
			ResultErr:   fmt.Errorf("[%s] assignment data is nil", Name),
			ErrData:     nil,
		},
		// implement setter
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockSetterData{
				err: err,
			},
			ResultErr: err,
			ErrData: &MockSetterData{
				err: err,
			},
			ExpectedData: &MockSetterData{
				err: err,
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockSetterData{},
			ErrData:     &MockSetterData{},
			ExpectedData: &MockSetterData{
				Key:   "key",
				Value: "value",
			},
		},
		// struct
		// struct err
		{
			Key:          "Key",
			Value:        "value",
			ParsedValue:  "value",
			Data:         MockData{},
			ResultErr:    fmt.Errorf("[%s] assignent not supported %v", Name, reflect.TypeOf(MockData{}).String()),
			ErrData:      MockData{},
			ExpectedData: MockData{},
		},
		{
			Key:         "Age",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] %v path value not exists key %s", Name, reflect.TypeOf(&MockData{}).String(), "Age"),
			ErrData:     &MockData{},
		},
		{
			Key:         "Age.Name",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] assignment %v not exists key %s", Name, reflect.TypeOf(&MockData{}).String(), "Age"),
			ErrData:     &MockData{},
		},
		{
			Key:         "Data.Name",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] %v path value %v can not set", Name, reflect.TypeOf(&MockData{}).String(), reflect.TypeOf(&MockData{}).String()),
			ErrData:     &MockData{},
		},
		{
			Key:         "NilData.Name",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] assignment data is nil", Name),
			ErrData:     &MockData{},
		},
		// struct success
		{
			Key:         "Key",
			Value:       "key",
			ParsedValue: "key",
			Data:        &MockData{},
			ExpectedData: &MockData{
				Key: "key",
			},
		},
		{
			Key:         "Value",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ExpectedData: &MockData{
				Value: "value",
			},
		},
		{
			Key:         "NilData.Name",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockData{
				NilData: &SubData{},
			},
			ExpectedData: &MockData{
				NilData: &SubData{
					Name: "value",
				},
			},
		},
		{
			Key:         "NilData.Data.Name",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockData{
				NilData: &SubData{
					Data: &SubData{},
				},
			},
			ExpectedData: &MockData{
				NilData: &SubData{
					Data: &SubData{
						Name: "value",
					},
				},
			},
		},
		{
			Key:         "nil_data.data.name",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockData{
				NilData: &SubData{
					Data: &SubData{},
				},
			},
			ExpectedData: &MockData{
				NilData: &SubData{
					Data: &SubData{
						Name: "value",
					},
				},
			},
		},
		{
			Key:         "nil_data.name",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockData{
				NilData: &SubData{},
			},
			ExpectedData: &MockData{
				NilData: &SubData{
					Name: "value",
				},
			},
		},
		// map
		{
			Key:         "key",
			Value:       "key",
			ParsedValue: "key",
			Data:        map[string]string{},
			ExpectedData: map[string]string{
				"key": "key",
			},
		},
		{
			Key:         "value",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]string{
				"key": "key",
			},
			ExpectedData: map[string]string{
				"key":   "key",
				"value": "value",
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]string{
				"key": "key",
			},
			ExpectedData: map[string]string{
				"key": "value",
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]string{
				"key": "key",
			},
			ExpectedData: map[string]string{
				"key": "value",
			},
		},

		{
			Key:         "key",
			Value:       "key",
			ParsedValue: "key",
			Data:        map[string]string{},
			ExpectedData: map[string]string{
				"key": "key",
			},
		},
		{
			Key:         "value",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]interface{}{
				"key": "key",
			},
			ExpectedData: map[string]interface{}{
				"key":   "key",
				"value": "value",
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]interface{}{
				"key": "key",
			},
			ExpectedData: map[string]interface{}{
				"key": "value",
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: map[string]interface{}{
				"key": "key",
			},
			ExpectedData: map[string]interface{}{
				"key": "value",
			},
		},
		{
			Key:         "key",
			Value:       1,
			ParsedValue: 1,
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": 1,
			},
		},
		{
			Key:         "key",
			Value:       1.0,
			ParsedValue: 1.0,
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": 1.0,
			},
		},
		{
			Key:         "key",
			Value:       -1,
			ParsedValue: -1,
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": -1,
			},
		},
		{
			Key:         "key",
			Value:       -1.0,
			ParsedValue: -1.0,
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": -1.0,
			},
		},
		{
			Key:         "key",
			Value:       []interface{}{"1"},
			ParsedValue: []interface{}{"1"},
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": []interface{}{"1"},
			},
		},
		{
			Key:         "key",
			Value:       []interface{}{1},
			ParsedValue: []interface{}{1},
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": []interface{}{1},
			},
		},
		{
			Key:         "key",
			Value:       []interface{}{1.0},
			ParsedValue: []interface{}{1.0},
			Data:        map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"key": []interface{}{1.0},
			},
		},
		// slice
		// err
		{
			Key:         "1",
			Value:       []string{},
			ParsedValue: []string{},
			ResultErr: fmt.Errorf("[%s] assignment %v path value %v length is %d but set index is %s",
				Name,
				reflect.TypeOf([]string{}).String(),
				reflect.TypeOf([]string{}).String(),
				0,
				"1"),
			Data:    []string{},
			ErrData: []string{},
		},
		{
			Key:         "h",
			Value:       []string{},
			ParsedValue: []string{},
			ResultErr: fmt.Errorf("[%s] assignment %v path value %v is a list but key [%s] can not convert to int",
				Name,
				reflect.TypeOf([]string{}).String(),
				reflect.TypeOf([]string{}).String(),
				"h"),
			Data:    []string{},
			ErrData: []string{},
		},
		{
			Key:          "0",
			Value:        "0",
			ParsedValue:  "0",
			Data:         make([]string, 1),
			ExpectedData: []string{"0"},
		},
		{
			Key:          "0",
			Value:        "0",
			ParsedValue:  "0",
			Data:         []string{"1"},
			ExpectedData: []string{"0"},
		},
		{
			Key:          "1",
			Value:        "1",
			ParsedValue:  "1",
			Data:         []string{"0", "0"},
			ExpectedData: []string{"0", "1"},
		},
	}

	for index, tt := range cases {
		parsedValue, err := ass.PrepareValue(ctx, tt.Value)
		assert.Nil(t, err)
		assert.Equal(t, tt.ParsedValue, tt.Value)

		err = ass.Run(ctx, tt.Data, tt.Key, parsedValue)
		if tt.ResultErr != nil {
			assert.True(t, reflect.DeepEqual(tt.ResultErr, err), index)
			assert.Equal(t, tt.ErrData, tt.Data, index)
		} else {
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(tt.Data, tt.ExpectedData), index)
		}
	}
}
