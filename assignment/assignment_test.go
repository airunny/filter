package assignment

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Work struct {
	WorkName string
}

type User struct {
	Name   string
	Age    int8
	IDCard string
	Works  []*Work
}

type Temp struct {
	User  *User
	Tests []string
}

func TestEqual_Run(t *testing.T) {
	structData := &Temp{
		User: &User{
			Name:   "zhangsan",
			Age:    18,
			IDCard: "110",
			Works: []*Work{
				{
					WorkName: "operation",
				},
			},
		},
		Tests: []string{"1", "2"},
	}

	mapData := map[string]interface{}{"user": structData, "name": "张三", "age": 18, "citys": []string{"1", "2"}, "ages": []int64{1, 2}}
	arrayData := []interface{}{structData, 1, "2", []int{1, 2}, []string{"1", "2"}}
	//arrayData1 := []string{"1", "2"}

	cases := []struct {
		Data     interface{}
		Key      string
		Value    interface{}
		Expected interface{}
		GotFunc  func() interface{}
	}{
		// map
		{
			Data:     mapData,
			Key:      "name",
			Value:    "李四",
			Expected: "李四",
			GotFunc: func() interface{} {
				return mapData["name"]
			},
		},
		{
			Data:     mapData,
			Key:      "age",
			Value:    19,
			Expected: 19,
			GotFunc: func() interface{} {
				return mapData["age"]
			},
		},
		{
			Data:     mapData,
			Key:      "citys",
			Value:    []string{"3", "4"},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return mapData["citys"]
			},
		},
		{
			Data:     mapData,
			Key:      "ages",
			Value:    []int64{3, 4},
			Expected: []int64{3, 4},
			GotFunc: func() interface{} {
				return mapData["ages"]
			},
		},
		{
			Data:     mapData,
			Key:      "user.User.Name",
			Value:    "zhangsan1",
			Expected: "zhangsan1",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data: mapData,
			Key:  "user.User",
			Value: &User{
				Name:   "user",
				Age:    19,
				IDCard: "120",
				Works: []*Work{
					{
						WorkName: "operation1",
					},
				},
			},
			Expected: &User{
				Name:   "user",
				Age:    19,
				IDCard: "120",
				Works: []*Work{
					{
						WorkName: "operation1",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     mapData,
			Key:      "user.Tests",
			Value:    []string{"3", "4"},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return structData.Tests
			},
		},
		{
			Data:     mapData,
			Key:      "user.User.Age",
			Value:    29,
			Expected: int8(29),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},

		// array
		{
			Data:     arrayData,
			Key:      "0.User.Name",
			Value:    "zhangsan2",
			Expected: "zhangsan2",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data: arrayData,
			Key:  "0.User",
			Value: &User{
				Name:   "user2",
				Age:    20,
				IDCard: "1200",
				Works: []*Work{
					{
						WorkName: "operation2",
					},
				},
			},
			Expected: &User{
				Name:   "user2",
				Age:    20,
				IDCard: "1200",
				Works: []*Work{
					{
						WorkName: "operation2",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     arrayData,
			Key:      "0.Tests",
			Value:    []string{"30", "40"},
			Expected: []string{"30", "40"},
			GotFunc: func() interface{} {
				return structData.Tests
			},
		},
		{
			Data:     arrayData,
			Key:      "0.User.Age",
			Value:    40,
			Expected: int8(40),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},
		{
			Data:     arrayData,
			Key:      "1",
			Value:    10,
			Expected: 10,
			GotFunc: func() interface{} {
				return arrayData[1]
			},
		},
		{
			Data:     arrayData,
			Key:      "2",
			Value:    "20",
			Expected: "20",
			GotFunc: func() interface{} {
				return arrayData[2]
			},
		},
		{
			Data:     arrayData,
			Key:      "3",
			Value:    []int{3, 4},
			Expected: []int{3, 4},
			GotFunc: func() interface{} {
				return arrayData[3]
			},
		},
		{
			Data:     arrayData,
			Key:      "4",
			Value:    []string{"3", "4"},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return arrayData[4]
			},
		},

		// ptr
		{
			Data:     structData,
			Key:      "User.Name",
			Value:    "golang",
			Expected: "golang",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data:     structData,
			Key:      "User.Age",
			Value:    17,
			Expected: int8(17),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},
		{
			Data:     structData,
			Key:      "User.IDCard",
			Value:    "id_card",
			Expected: "id_card",
			GotFunc: func() interface{} {
				return structData.User.IDCard
			},
		},
		{
			Data: structData,
			Key:  "User",
			Value: &User{
				Name:   "golang1",
				Age:    10,
				IDCard: "2009",
				Works: []*Work{
					{
						WorkName: "shanghai",
					},
				},
			},
			Expected: &User{
				Name:   "golang1",
				Age:    10,
				IDCard: "2009",
				Works: []*Work{
					{
						WorkName: "shanghai",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     structData,
			Key:      "User.Works.0.WorkName",
			Value:    "suzhou",
			Expected: "suzhou",
			GotFunc: func() interface{} {
				return structData.User.Works[0].WorkName
			},
		},
	}

	instance := &Equal{}
	for index, v := range cases {
		prepayValue, err := instance.PrepareValue(context.Background(), v.Value)
		assert.Equal(t, nil, err, index)
		instance.Run(context.Background(), v.Data, v.Key, prepayValue)
		if !reflect.DeepEqual(v.Expected, v.GotFunc()) {
			t.Errorf("%v expected %v but got %v", index, v.Expected, v.GotFunc())
		}
	}
}
