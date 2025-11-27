package executor

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/liyanbing/filter/assignment/set"
	"github.com/stretchr/testify/assert"
)

type MockData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGroup(t *testing.T) {
	cases := []struct {
		BuildErr     error
		ExecuteErr   error
		Items        []interface{}
		Data         interface{}
		ExpectedData interface{}
	}{
		// err
		{
			BuildErr: fmt.Errorf("executor group item must be array"),
			Items: []interface{}{
				"name",
				"=",
				"name",
			},
		},
		// struct success
		{
			Items: []interface{}{
				[]interface{}{"name", "=", "name"},
				[]interface{}{"age", "=", 10},
			},
			Data: &MockData{},
			ExpectedData: &MockData{
				Name: "name",
				Age:  10,
			},
		},
		{
			Items: []interface{}{
				[]interface{}{"name", "=", "golang"},
				[]interface{}{"age", "=", 100},
			},
			Data: &MockData{},
			ExpectedData: &MockData{
				Name: "golang",
				Age:  100,
			},
		},
		{
			Items: []interface{}{
				[]interface{}{"name", "=", "name"},
			},
			Data: &MockData{},
			ExpectedData: &MockData{
				Name: "name",
			},
		},
		{
			Items: []interface{}{
				[]interface{}{"age", "=", 10},
			},
			Data: &MockData{},
			ExpectedData: &MockData{
				Age: 10,
			},
		},
		// map success
		{
			Items: []interface{}{
				[]interface{}{"age", "=", 10},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age": 10,
			},
		},
		{
			Items: []interface{}{
				[]interface{}{"name", "=", "golang"},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "golang",
			},
		},
		{
			Items: []interface{}{
				[]interface{}{"name", "=", "golang"},
				[]interface{}{"age", "=", 10},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "golang",
				"age":  10,
			},
		},
	}

	ctx := context.Background()
	for index, tt := range cases {
		execute, err := BuildGroup(ctx, tt.Items)
		if err != nil {
			if index == 1 {
				fmt.Println(tt.BuildErr)
				fmt.Println(err)
			}
			assert.True(t, reflect.DeepEqual(tt.BuildErr, err), index)
		} else {
			err = execute.Execute(ctx, tt.Data)
			assert.True(t, reflect.DeepEqual(tt.ExecuteErr, err), index)
			assert.True(t, reflect.DeepEqual(tt.ExpectedData, tt.Data), index)
		}
	}
}
