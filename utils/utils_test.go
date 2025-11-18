package utils

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatEquals(t *testing.T) {
	cases := []struct {
		A        float64
		B        float64
		Expected bool
	}{
		{
			A:        0,
			B:        0 + EPSILON,
			Expected: false,
		},
		{
			A:        0,
			B:        0 + 0.000000001,
			Expected: true,
		},
		{
			A:        -1,
			B:        -1 + 0.000000001,
			Expected: true,
		},
		{
			A:        math.MaxFloat64,
			B:        math.MaxFloat64 - 0.0000001,
			Expected: true,
		},
		{
			A:        1,
			B:        1,
			Expected: true,
		},
	}

	for i, v := range cases {
		assert.Equal(t, v.Expected, FloatEquals(v.A, v.B), i)
	}
}

type WeightMock struct {
	Weight int64
}

func (s WeightMock) GetWeight() int64 {
	return s.Weight
}

func TestTotalWeight(t *testing.T) {
	weightArr := make([]IWeight, 0, 10)
	total := int64(0)
	for i := 0; i < 10; i++ {
		total += int64(i)
		weightArr = append(weightArr, &WeightMock{
			Weight: int64(i),
		})
	}

	assert.Equal(t, total, TotalWeight(weightArr))
}

func TestPickByWeight(t *testing.T) {
	weightArr := make([]IWeight, 0, 10)
	total := int64(0)
	for i := 1; i <= 10; i++ {
		total += int64(i)
		weightArr = append(weightArr, &WeightMock{
			Weight: int64(i),
		})
	}

	pickCache := make(map[int]int)
	totalCount := 100000
	for i := 0; i < totalCount; i++ {
		pickIndex := PickByWeight(weightArr, total)
		pickIndex++
		if _, ok := pickCache[pickIndex]; ok {
			pickCache[pickIndex]++
		} else {
			pickCache[pickIndex] = 1
		}
	}

	for k, v := range pickCache {
		expected := float64(k) / float64(total)
		got := float64(v) / float64(totalCount)
		assert.Equal(t, true, (expected-got) < 1 && (got-expected) < 1)
	}
}

// -------------
type Work struct {
	Name string `json:"name"`
}

type User struct {
	Name   string `json:"name"`
	Age    int32  `json:"age"`
	IDCard string `json:"id_card"`
	Works  []Work `json:"works"`
}

type Temp struct {
	User User `json:"user"`
}

func TestGetObjectValueByKey(t *testing.T) {
	mock := Temp{
		User: User{
			Name:   "zhangsan",
			Age:    18,
			IDCard: "110",
			Works: []Work{
				{
					Name: "111",
				},
				{
					Name: "222",
				},
			},
		}}

	cases := []struct {
		Data     interface{}
		Key      string
		Expected interface{}
		OK       bool
	}{
		// map
		{
			Data:     map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "user",
			Expected: "zhangsan",
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "age",
			Expected: 18,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// map ptr
		{
			Data:     &map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "user",
			Expected: "zhangsan",
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "age",
			Expected: 18,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// array
		{
			Data:     []interface{}{mock},
			Key:      "0",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// array ptr
		{
			Data:     &[]interface{}{mock},
			Key:      "0",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// struct
		{
			Data:     mock,
			Key:      ".",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// struct ptr
		{
			Data:     &mock,
			Key:      ".",
			Expected: &mock,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "",
			Expected: &mock,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
	}

	for index, v := range cases {
		ret, ok := GetObjectValueByKey(context.Background(), v.Data, v.Key)
		assert.Equal(t, v.OK, ok, index)
		assert.Equal(t, true, reflect.DeepEqual(v.Expected, ret))
	}
}
