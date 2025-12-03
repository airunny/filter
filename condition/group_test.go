package condition

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/operations/equal"
	_ "github.com/airunny/filter/operations/greater_than"
	_ "github.com/airunny/filter/operations/less_than"
	_ "github.com/airunny/filter/variables/success"
	_ "github.com/airunny/filter/variables/time"
	"github.com/stretchr/testify/assert"
)

func TestBuildGroup(t *testing.T) {
	cases := []struct {
		Items  []interface{}
		Result bool
		Logic  Logic
		Err    error
	}{
		// err
		{
			Items: []interface{}{
				"1",
				"2",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    errors.New("condition item is not array"),
		},
		// and
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", ">", 1},
			},
			Result: true,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", ">", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		// or
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", ">", 1},
			},
			Result: true,
			Logic:  LogicOr,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: true,
			Logic:  LogicOr,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", ">", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: false,
			Logic:  LogicOr,
		},
		// not
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", ">", 1},
			},
			Result: false,
			Logic:  LogicNot,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: false,
			Logic:  LogicNot,
		},
		{
			Items: []interface{}{
				[]interface{}{"success", ">", 1},
				[]interface{}{"timestamp", "<", 1},
			},
			Result: true,
			Logic:  LogicNot,
		},
	}

	ctx := context.Background()
	cc := cache.NewCache()
	for index, tt := range cases {
		cond, err := BuildGroup(ctx, tt.Items, tt.Logic)
		if err != nil {
			assert.True(t, reflect.DeepEqual(tt.Err, err))
		} else {
			ok, err := cond.IsConditionOk(ctx, nil, cc)
			assert.Nil(t, err)
			assert.Equal(t, tt.Result, ok, index)
		}
	}
}
