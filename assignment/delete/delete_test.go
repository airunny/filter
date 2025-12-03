package delete

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/airunny/filter/assignment"
	"github.com/stretchr/testify/assert"
)

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

type MockDeleterData struct {
	Name string
	Err  error
}

func (m *MockDeleterData) Delete(key string, value interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	m.Name = ""
	return nil
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
		{
			Key:         "NilData.Age",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] assignment data is nil", Name),
			ErrData:     &MockData{},
		},
		{
			Key:         "NullData.Age",
			Value:       "value",
			ParsedValue: "value",
			Data:        &MockData{},
			ResultErr:   fmt.Errorf("[%s] assignment %v not exists key %s", Name, reflect.TypeOf(&MockData{}).String(), "NullData"),
			ErrData:     &MockData{},
		},
		// implement setter
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockDeleterData{
				Err: err,
			},
			ResultErr: err,
			ErrData: &MockDeleterData{
				Err: err,
			},
		},
		{
			Key:         "key",
			Value:       "value",
			ParsedValue: "value",
			Data: &MockDeleterData{
				Name: "name",
			},
			ExpectedData: &MockDeleterData{
				Name: "",
			},
		},
		// map
		{
			Key:         "key",
			Value:       "key",
			ParsedValue: "key",
			Data: map[string]interface{}{
				"key": "key",
			},
			ExpectedData: map[string]interface{}{},
		},
		{
			Key:         "value",
			Value:       "key",
			ParsedValue: "key",
			Data: map[string]interface{}{
				"key": "key",
			},
			ExpectedData: map[string]interface{}{
				"key": "key",
			},
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
