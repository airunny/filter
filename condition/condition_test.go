package condition

import (
	"context"
	"errors"
	"fmt"
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

func TestBuildCondition(t *testing.T) {
	cases := []struct {
		Items  []interface{}
		Result bool
		Logic  Logic
		Err    error
	}{
		// err
		{
			Items:  []interface{}{},
			Result: true,
			Logic:  LogicAnd,
			Err:    errors.New("condition is empty"),
		},
		{
			Items: []interface{}{
				"success",
				"=",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    errors.New("condition item must contains three element"),
		},
		{
			Items: []interface{}{
				1,
				"=",
				1,
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("condition item 1st element[%v] is not string", 1),
		},
		{
			Items: []interface{}{
				"and",
				"=",
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("group condition [%s] 3rd element is not array", "and"),
		},
		{
			Items: []interface{}{
				"or",
				"=",
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("group condition [%s] 3rd element is not array", "or"),
		},
		{
			Items: []interface{}{
				"not",
				"=",
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("group condition [%s] 3rd element is not array", "not"),
		},
		{
			Items: []interface{}{
				"golang",
				"=",
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("condition not exists variable [%s]", "golang"),
		},
		{
			Items: []interface{}{
				"time",
				1,
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("condition operation should be string [%v]", 1),
		},
		{
			Items: []interface{}{
				"time",
				"match",
				"1",
			},
			Result: true,
			Logic:  LogicAnd,
			Err:    fmt.Errorf("condition not exists operation [%s]", "match"),
		},
		// condition
		{
			Items: []interface{}{
				"success",
				"=",
				1,
			},
			Result: true,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"success",
				">",
				1,
			},
			Result: false,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"timestamp",
				">",
				1,
			},
			Result: true,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"timestamp",
				"<",
				1,
			},
			Result: false,
			Logic:  LogicAnd,
		},
		// group
		// and
		{
			Items: []interface{}{
				"and",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", ">", 1},
				},
			},
			Result: true,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"and",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"and",
				"=>",
				[]interface{}{
					[]interface{}{"success", ">", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		// or
		{
			Items: []interface{}{
				"or",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", ">", 1},
				},
			},
			Result: true,
			Logic:  LogicOr,
		},
		{
			Items: []interface{}{
				"or",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: true,
			Logic:  LogicOr,
		},
		{
			Items: []interface{}{
				"or",
				"=>",
				[]interface{}{
					[]interface{}{"success", ">", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: false,
			Logic:  LogicOr,
		},
		// not
		{
			Items: []interface{}{
				"not",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", ">", 1},
				},
			},
			Result: false,
			Logic:  LogicNot,
		},
		{
			Items: []interface{}{
				"not",
				"=>",
				[]interface{}{
					[]interface{}{"success", "=", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: false,
			Logic:  LogicNot,
		},
		{
			Items: []interface{}{
				"not",
				"=>",
				[]interface{}{
					[]interface{}{"success", ">", 1},
					[]interface{}{"timestamp", "<", 1},
				},
			},
			Result: true,
			Logic:  LogicNot,
		},
		{
			Items: []interface{}{
				"or",
				"=>",
				[]interface{}{
					[]interface{}{ // true
						"and",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
						},
					},
					[]interface{}{ // false
						"not",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<", 1},
						},
					},
				},
			},
			Result: true,
			Logic:  LogicOr,
		},
		{
			Items: []interface{}{
				"and",
				"=>",
				[]interface{}{
					[]interface{}{ // true
						"and",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
						},
					},
					[]interface{}{ // false
						"not",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<", 1},
						},
					},
				},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"not",
				"=>",
				[]interface{}{
					[]interface{}{ // true
						"and",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
						},
					},
					[]interface{}{ // false
						"not",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<", 1},
						},
					},
				},
			},
			Result: false,
			Logic:  LogicAnd,
		},
		{
			Items: []interface{}{
				"not",
				"=>",
				[]interface{}{
					[]interface{}{ // false
						"and",
						"=>",
						[]interface{}{
							[]interface{}{"success", "<", 1},
							[]interface{}{"timestamp", "<", 1},
						},
					},
					[]interface{}{ // false
						"not",
						"=>",
						[]interface{}{
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<", 1},
						},
					},
				},
			},
			Result: true,
			Logic:  LogicNot,
		},
	}

	ctx := context.Background()
	cc := cache.NewCache()
	for index, tt := range cases {
		cond, err := BuildCondition(ctx, tt.Items, tt.Logic)
		if err != nil {
			assert.True(t, reflect.DeepEqual(tt.Err, err), index)
		} else {
			ok, err := cond.IsConditionOk(ctx, nil, cc)
			assert.Nil(t, err)
			assert.Equal(t, tt.Result, ok, index)
		}
	}
}
